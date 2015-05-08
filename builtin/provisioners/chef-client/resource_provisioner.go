package chefclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/terraform/communicator"
	"github.com/hashicorp/terraform/communicator/remote"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/go-linereader"
	"github.com/mitchellh/mapstructure"
)

const (
	firstBoot      = "first-boot.json"
	logfileDir     = "logfiles"
	linuxConfDir   = "/etc/chef"
	windowsConfDir = "C:/chef"
)

const clientConf = `
log_location            STDOUT
chef_server_url         "{{ .ServerURL }}"
validation_client_name  "{{ .ValidationClientName }}"
node_name               "{{ .NodeName }}"

{{ if .HTTPProxy }}
http_proxy          "{{ .HTTPProxy }}"
ENV['http_proxy'] = "{{ .HTTPProxy }}"
ENV['HTTP_PROXY'] = "{{ .HTTPProxy }}"
{{ end }}

{{ if .HTTPSProxy }}
https_proxy          "{{ .HTTPSProxy }}"
ENV['https_proxy'] = "{{ .HTTPSProxy }}"
ENV['HTTPS_PROXY'] = "{{ .HTTPSProxy }}"
{{ end }}

{{ if .NOProxy }}no_proxy "{{ join .NOProxy "," }}"{{ end }}
{{ if .SSLVerifyMode }}ssl_verify_mode {{ .SSLVerifyMode }}{{ end }}
`

// Provisioner represents a specificly configured chef provisioner
type Provisioner struct {
	Attributes           interface{} `mapstructure:"-"`
	Environment          string      `mapstructure:"environment"`
	LogToFile            bool        `mapstructure:"log_to_file"`
	HTTPProxy            string      `mapstructure:"http_proxy"`
	HTTPSProxy           string      `mapstructure:"https_proxy"`
	NOProxy              []string    `mapstructure:"no_proxy"`
	NodeName             string      `mapstructure:"node_name"`
	PreventSudo          bool        `mapstructure:"prevent_sudo"`
	RunList              []string    `mapstructure:"run_list"`
	ServerURL            string      `mapstructure:"server_url"`
	SkipInstall          bool        `mapstructure:"skip_install"`
	SSLVerifyMode        string      `mapstructure:"ssl_verify_mode"`
	ValidationClientName string      `mapstructure:"validation_client_name"`
	ValidationKeyPath    string      `mapstructure:"validation_key_path"`
	Version              string      `mapstructure:"version"`

	installChefClient func(terraform.UIOutput, communicator.Communicator) error
	createConfigFiles func(terraform.UIOutput, communicator.Communicator) error
	runChefClient     func(terraform.UIOutput, communicator.Communicator) error
}

// ResourceProvisioner represents a generic chef provisioner
type ResourceProvisioner struct{}

// Apply executes the file provisioner
func (r *ResourceProvisioner) Apply(
	o terraform.UIOutput,
	s *terraform.InstanceState,
	c *terraform.ResourceConfig) error {
	// Decode the raw config for this provisioner
	p, err := r.decodeConfig(c)
	if err != nil {
		return err
	}

	// Set some values based on the targeted OS
	switch s.Ephemeral.ConnInfo["type"] {
	case "ssh", "": // The default connection type is ssh, so if the type is empty use ssh
		p.PreventSudo = p.PreventSudo || s.Ephemeral.ConnInfo["user"] == "root"
		p.installChefClient = p.sshInstallChefClient
		p.createConfigFiles = p.sshCreateConfigFiles
		p.runChefClient = p.runChefClientFunc(linuxConfDir)
	case "winrm":
		p.PreventSudo = true
		p.installChefClient = p.winrmInstallChefClient
		p.createConfigFiles = p.winrmCreateConfigFiles
		p.runChefClient = p.runChefClientFunc(windowsConfDir)
	default:
		return fmt.Errorf("Unsupported connection type: %s", s.Ephemeral.ConnInfo["type"])
	}

	// Get a new communicator
	comm, err := communicator.New(s)
	if err != nil {
		return err
	}

	// Wait and retry until we establish the connection
	err = retryFunc(comm.Timeout(), func() error {
		err := comm.Connect(o)
		return err
	})
	if err != nil {
		return err
	}
	defer comm.Disconnect()

	if !p.SkipInstall {
		if err := p.installChefClient(o, comm); err != nil {
			return err
		}
	}

	o.Output("Creating configuration files...")
	if err := p.createConfigFiles(o, comm); err != nil {
		return err
	}

	o.Output("Starting initial Chef-Client run...")
	if err := p.runChefClient(o, comm); err != nil {
		return err
	}

	return nil
}

// Validate checks if the required arguments are configured
func (r *ResourceProvisioner) Validate(c *terraform.ResourceConfig) (ws []string, es []error) {
	p, err := r.decodeConfig(c)
	if err != nil {
		es = append(es, err)
		return ws, es
	}

	if p.NodeName == "" {
		es = append(es, fmt.Errorf("Key not found: node_name"))
	}
	if p.RunList == nil {
		es = append(es, fmt.Errorf("Key not found: run_list"))
	}
	if p.ServerURL == "" {
		es = append(es, fmt.Errorf("Key not found: server_url"))
	}
	if p.ValidationClientName == "" {
		es = append(es, fmt.Errorf("Key not found: validation_client_name"))
	}
	if p.ValidationKeyPath == "" {
		es = append(es, fmt.Errorf("Key not found: validation_key_path"))
	}

	return ws, es
}

func (r *ResourceProvisioner) decodeConfig(c *terraform.ResourceConfig) (*Provisioner, error) {
	p := new(Provisioner)

	decConf := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           p,
	}
	dec, err := mapstructure.NewDecoder(decConf)
	if err != nil {
		return nil, err
	}

	if err := dec.Decode(c.Raw); err != nil {
		return nil, err
	}

	if p.Environment == "" {
		p.Environment = "_default"
	}

	if attrs, ok := c.Raw["attributes"]; ok {
		p.Attributes, err = rawToJSON(attrs)
		if err != nil {
			return nil, fmt.Errorf("Error parsing the attributes: %v", err)
		}
	}

	return p, nil
}

func rawToJSON(raw interface{}) (interface{}, error) {
	switch s := raw.(type) {
	case []map[string]interface{}:
		if len(s) != 1 {
			return nil, errors.New("unexpected input while parsing raw config to JSON")
		}

		var err error
		for k, v := range s[0] {
			s[0][k], err = rawToJSON(v)
			if err != nil {
				return nil, err
			}
		}

		return s[0], nil
	default:
		return raw, nil
	}
}

// retryFunc is used to retry a function for a given duration
func retryFunc(timeout time.Duration, f func() error) error {
	finish := time.After(timeout)
	for {
		err := f()
		if err == nil {
			return nil
		}
		log.Printf("Retryable error: %v", err)

		select {
		case <-finish:
			return err
		case <-time.After(3 * time.Second):
		}
	}
}

func (p *Provisioner) runChefClientFunc(
	confDir string) func(terraform.UIOutput, communicator.Communicator) error {
	return func(o terraform.UIOutput, comm communicator.Communicator) error {
		fb := path.Join(confDir, firstBoot)
		cmd := fmt.Sprintf("chef-client -j %q -E %q", fb, p.Environment)

		if p.LogToFile {
			if err := os.MkdirAll(logfileDir, 0777); err != nil {
				return fmt.Errorf("Error creating logfile directory %s: %v", logfileDir, err)
			}

			logFile := path.Join(logfileDir, p.NodeName)
			f, err := os.Create(path.Join(logFile))
			if err != nil {
				return fmt.Errorf("Error creating logfile %s: %v", logFile, err)
			}
			f.Close()

			o.Output("Writing Chef Client output to " + logFile)
			o = p
		}

		return p.runCommand(o, comm, cmd)
	}
}

// Output implementation of terraform.UIOutput interface
func (p *Provisioner) Output(output string) {
	logFile := path.Join(logfileDir, p.NodeName)
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Error creating logfile %s: %v", logFile, err)
		return
	}
	defer f.Close()

	// These steps are needed to remove any ANSI escape codes used to colorize
	// the output and to make sure we have proper line endings before writing
	// the string to the logfile.
	re := regexp.MustCompile(`\x1b\[[0-9;]+m`)
	output = re.ReplaceAllString(output, "")
	output = strings.Replace(output, "\r", "\n", -1)

	if _, err := f.WriteString(output); err != nil {
		log.Printf("Error writing output to logfile %s: %v", logFile, err)
	}

	if err := f.Sync(); err != nil {
		log.Printf("Error saving logfile %s to disk: %v", logFile, err)
	}
}

func (p *Provisioner) deployConfigFiles(
	o terraform.UIOutput,
	comm communicator.Communicator,
	confDir string) error {
	// Open the validation .pem file
	f, err := os.Open(p.ValidationKeyPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Copy the validation .pem to the new instance
	if err := comm.Upload(path.Join(confDir, "validation.pem"), f); err != nil {
		return fmt.Errorf("Uploading validation.pem failed: %v", err)
	}

	// Make strings.Join available for use within the template
	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	// Create a new template and parse the client.rb into it
	t := template.Must(template.New("client.rb").Funcs(funcMap).Parse(clientConf))

	var buf bytes.Buffer
	err = t.Execute(&buf, p)
	if err != nil {
		return fmt.Errorf("Error executing client.rb template: %s", err)
	}

	// Copy the client.rb to the new instance
	if err := comm.Upload(path.Join(confDir, "client.rb"), &buf); err != nil {
		return fmt.Errorf("Uploading client.rb failed: %v", err)
	}

	// Create a map with first boot settings
	fb := make(map[string]interface{})
	if p.Attributes != nil {
		fb = p.Attributes.(map[string]interface{})
	}

	// Add the initial runlist to the first boot settings
	fb["run_list"] = p.RunList

	// Marshal the first boot settings to JSON
	d, err := json.Marshal(fb)
	if err != nil {
		return fmt.Errorf("Failed to create %s data: %s", firstBoot, err)
	}

	// Copy the first-boot.json to the new instance
	if err := comm.Upload(path.Join(confDir, firstBoot), bytes.NewReader(d)); err != nil {
		return fmt.Errorf("Uploading first-boot.json failed: %v", err)
	}

	return nil
}

// runCommand is used to run already prepared commands
func (p *Provisioner) runCommand(
	o terraform.UIOutput,
	comm communicator.Communicator,
	command string) error {
	var err error

	// Unless prevented, prefix the command with sudo
	if !p.PreventSudo {
		command = "sudo " + command
	}

	outR, outW := io.Pipe()
	errR, errW := io.Pipe()
	outDoneCh := make(chan struct{})
	errDoneCh := make(chan struct{})
	go p.copyOutput(o, outR, outDoneCh)
	go p.copyOutput(o, errR, errDoneCh)

	cmd := &remote.Cmd{
		Command: command,
		Stdout:  outW,
		Stderr:  errW,
	}

	if err := comm.Start(cmd); err != nil {
		return fmt.Errorf("Error executing command %q: %v", cmd.Command, err)
	}

	cmd.Wait()
	if cmd.ExitStatus != 0 {
		err = fmt.Errorf(
			"Command %q exited with non-zero exit status: %d", cmd.Command, cmd.ExitStatus)
	}

	// Wait for output to clean up
	outW.Close()
	errW.Close()
	<-outDoneCh
	<-errDoneCh

	// If we have an error, return it out now that we've cleaned up
	if err != nil {
		return err
	}

	return nil
}

func (p *Provisioner) copyOutput(o terraform.UIOutput, r io.Reader, doneCh chan<- struct{}) {
	defer close(doneCh)
	lr := linereader.New(r)
	for line := range lr.Ch {
		o.Output(line)
	}
}
