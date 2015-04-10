package winrm

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hashicorp/terraform/communicator/remote"
	"github.com/hashicorp/terraform/terraform"
	"github.com/masterzen/winrm/winrm"
	"github.com/packer-community/winrmcp/winrmcp"
)

// Communicator represents the WinRM communicator
type Communicator struct {
	connInfo *connectionInfo
	client   *winrm.Client
	endpoint *winrm.Endpoint
}

// New creates a new communicator implementation over WinRM.
//func New(endpoint *winrm.Endpoint, user string, password string, timeout time.Duration) (*communicator, error) {
func New(s *terraform.InstanceState) (*Communicator, error) {
	connInfo, err := parseConnectionInfo(s)
	if err != nil {
		return nil, err
	}

	endpoint := &winrm.Endpoint{
		Host:     connInfo.Host,
		Port:     connInfo.Port,
		HTTPS:    connInfo.HTTPS,
		Insecure: connInfo.Insecure,
		CACert:   connInfo.CACert,
	}

	comm := &Communicator{
		connInfo: connInfo,
		endpoint: endpoint,
	}

	return comm, nil
}

// Connect implementation of communicator.Communicator interface
func (c *Communicator) Connect(o terraform.UIOutput) error {
	if c.client != nil {
		return nil
	}

	params := winrm.DefaultParameters()
	params.Timeout = formatDuration(c.Timeout())

	client, err := winrm.NewClientWithParameters(
		c.endpoint, c.connInfo.User, c.connInfo.Password, params)
	if err != nil {
		return err
	}

	if o != nil {
		o.Output(fmt.Sprintf(
			"Connecting to remote host via WinRM...\n"+
				"  Host: %s\n"+
				"  Port: %d\n"+
				"  User: %s\n"+
				"  Password: %t\n"+
				"  HTTPS: %t\n"+
				"  Insecure: %t\n"+
				"  CACert: %t",
			c.connInfo.Host,
			c.connInfo.Port,
			c.connInfo.User,
			c.connInfo.Password != "",
			c.connInfo.HTTPS,
			c.connInfo.Insecure,
			c.connInfo.CACert != nil,
		))
	}

	log.Printf("connecting to remote shell using WinRM")
	shell, err := client.CreateShell()
	if err != nil {
		log.Printf("connection error: %s", err)
		return err
	}

	err = shell.Close()
	if err != nil {
		log.Printf("error closing connection: %s", err)
		return err
	}

	if o != nil {
		o.Output("Connected!")
	}

	c.client = client

	return nil
}

// Disconnect implementation of communicator.Communicator interface
func (c *Communicator) Disconnect() error {
	c.client = nil
	return nil
}

// Timeout implementation of communicator.Communicator interface
func (c *Communicator) Timeout() time.Duration {
	return c.connInfo.TimeoutVal
}

// ScriptPath implementation of communicator.Communicator interface
func (c *Communicator) ScriptPath() string {
	return c.connInfo.ScriptPath
}

// Start implementation of communicator.Communicator interface
func (c *Communicator) Start(rc *remote.Cmd) error {
	log.Printf("starting remote command: %s", rc.Command)

	err := c.Connect(nil)
	if err != nil {
		return err
	}

	shell, err := c.client.CreateShell()
	if err != nil {
		return err
	}

	cmd, err := shell.Execute(rc.Command)
	if err != nil {
		return err
	}

	go runCommand(shell, cmd, rc)
	return nil
}

func runCommand(shell *winrm.Shell, cmd *winrm.Command, rc *remote.Cmd) {
	defer shell.Close()

	go io.Copy(rc.Stdout, cmd.Stdout)
	go io.Copy(rc.Stderr, cmd.Stderr)

	cmd.Wait()
	rc.SetExited(cmd.ExitCode())
}

// Upload implementation of communicator.Communicator interface
func (c *Communicator) Upload(path string, input io.Reader) error {
	wcp, err := c.newCopyClient()
	if err != nil {
		return err
	}
	return wcp.Write(path, input)
}

// UploadScript implementation of communicator.Communicator interface
func (c *Communicator) UploadScript(path string, input io.Reader) error {
	return c.Upload(path, input)
}

// UploadDir implementation of communicator.Communicator interface
func (c *Communicator) UploadDir(dst string, src string) error {
	log.Printf("Upload dir '%s' to '%s'", src, dst)
	wcp, err := c.newCopyClient()
	if err != nil {
		return err
	}
	return wcp.Copy(src, dst)
}

func (c *Communicator) newCopyClient() (*winrmcp.Winrmcp, error) {
	addr := fmt.Sprintf("%s:%d", c.endpoint.Host, c.endpoint.Port)
	return winrmcp.New(addr, &winrmcp.Config{
		Auth: winrmcp.Auth{
			User:     c.connInfo.User,
			Password: c.connInfo.Password,
		},
		OperationTimeout:      c.Timeout(),
		MaxOperationsPerShell: 15, // lowest common denominator
	})
}
