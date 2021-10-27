package logging

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/hashicorp/go-hclog"
)

// This output is shown if a panic happens.
const panicOutput = `
!!!!!!!!!!!!!!!!!!!!!!!!!!! TERRAFORM CRASH !!!!!!!!!!!!!!!!!!!!!!!!!!!!

Terraform crashed! This is always indicative of a bug within Terraform.
Please report the crash with Terraform[1] so that we can fix this.

When reporting bugs, please include your terraform version, the stack trace
shown below, and any additional information which may help replicate the issue.

[1]: https://github.com/hashicorp/terraform/issues

!!!!!!!!!!!!!!!!!!!!!!!!!!! TERRAFORM CRASH !!!!!!!!!!!!!!!!!!!!!!!!!!!!

`

// PanicHandler is called to recover from an internal panic in Terraform, and
// augments the standard stack trace with a more user friendly error message.
// PanicHandler must be called as a defered function.
func PanicHandler() {
	recovered := recover()
	if recovered == nil {
		return
	}

	fmt.Fprint(os.Stderr, panicOutput)
	fmt.Fprint(os.Stderr, recovered)

	// When called from a deferred function, debug.PrintStack will include the
	// full stack from the point of the pending panic.
	debug.PrintStack()
	os.Exit(2)
}

const pluginPanicOutput = `
Stack trace from the %[1]s plugin:

%s

Error: The %[1]s plugin crashed!

This is always indicative of a bug within the plugin. It would be immensely
helpful if you could report the crash with the plugin's maintainers so that it
can be fixed. The output above should help diagnose the issue.
`

// PluginPanics returns a series of provider panics that were collected during
// execution, and formatted for output.
func PluginPanics() []string {
	return panics.allPanics()
}

// panicRecorder provides a registry to check for plugin panics that may have
// happened when a plugin suddenly terminates.
type panicRecorder struct {
	sync.Mutex

	// panics maps the plugin name to the panic output lines received from
	// the logger.
	panics map[string][]string

	// maxLines is the max number of lines we'll record after seeing a
	// panic header. Since this is going to be printed in the UI output, we
	// don't want to destroy the scrollback. In most cases, the first few lines
	// of the stack trace is all that are required.
	maxLines int
}

// registerPlugin returns an accumulator function which will accept lines of
// a panic stack trace to collect into an error when requested.
func (p *panicRecorder) registerPlugin(name string) func(string) {
	p.Lock()
	defer p.Unlock()

	// In most cases we shouldn't be starting a plugin if it already
	// panicked, but clear out previous entries just in case.
	delete(p.panics, name)

	count := 0

	// this callback is used by the logger to store panic output
	return func(line string) {
		p.Lock()
		defer p.Unlock()

		// stop recording if there are too many lines.
		if count > p.maxLines {
			return
		}
		count++

		p.panics[name] = append(p.panics[name], line)
	}
}

func (p *panicRecorder) allPanics() []string {
	p.Lock()
	defer p.Unlock()

	var res []string
	for name, lines := range p.panics {
		if len(lines) == 0 {
			continue
		}

		res = append(res, fmt.Sprintf(pluginPanicOutput, name, strings.Join(lines, "\n")))
	}
	return res
}

// logPanicWrapper wraps an hclog.Logger and intercepts and records any output
// that appears to be a panic.
type logPanicWrapper struct {
	hclog.Logger
	panicRecorder func(string)
	inPanic       bool
}

// go-plugin will create a new named logger for each plugin binary.
func (l *logPanicWrapper) Named(name string) hclog.Logger {
	return &logPanicWrapper{
		Logger:        l.Logger.Named(name),
		panicRecorder: panics.registerPlugin(name),
	}
}

// we only need to implement Debug, since that is the default output level used
// by go-plugin when encountering unstructured output on stderr.
func (l *logPanicWrapper) Debug(msg string, args ...interface{}) {
	// We don't have access to the binary itself, so guess based on the stderr
	// output if this is the start of the traceback. An occasional false
	// positive shouldn't be a big deal, since this is only retrieved after an
	// error of some sort.

	panicPrefix := strings.HasPrefix(msg, "panic: ") || strings.HasPrefix(msg, "fatal error: ")

	l.inPanic = l.inPanic || panicPrefix

	if l.inPanic && l.panicRecorder != nil {
		l.panicRecorder(msg)
	}

	l.Logger.Debug(msg, args...)
}
