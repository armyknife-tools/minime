package main

import (
	"os"
	"os/signal"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/hashicorp/terraform/command"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available Terraform commands.
var Commands map[string]cli.CommandFactory

// Ui is the cli.Ui used for communicating to the outside world.
var Ui cli.Ui

const ErrorPrefix = "e:"
const OutputPrefix = "o:"

func init() {
	Ui = &cli.PrefixedUi{
		AskPrefix:    OutputPrefix,
		OutputPrefix: OutputPrefix,
		InfoPrefix:   OutputPrefix,
		ErrorPrefix:  ErrorPrefix,
		Ui:           &cli.BasicUi{Writer: os.Stdout},
	}

	meta := command.Meta{
		Color:       terminal.IsTerminal(int(os.Stdout.Fd())),
		ContextOpts: &ContextOpts,
		Ui:          Ui,
	}

	Commands = map[string]cli.CommandFactory{
		"apply": func() (cli.Command, error) {
			return &command.ApplyCommand{
				Meta:       meta,
				ShutdownCh: makeShutdownCh(),
			}, nil
		},

		"graph": func() (cli.Command, error) {
			return &command.GraphCommand{
				Meta: meta,
			}, nil
		},

		"output": func() (cli.Command, error) {
			return &command.OutputCommand{
				Meta: meta,
			}, nil
		},

		"plan": func() (cli.Command, error) {
			return &command.PlanCommand{
				Meta: meta,
			}, nil
		},

		"refresh": func() (cli.Command, error) {
			return &command.RefreshCommand{
				Meta: meta,
			}, nil
		},

		"show": func() (cli.Command, error) {
			return &command.ShowCommand{
				Meta: meta,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:              meta,
				Revision:          GitCommit,
				Version:           Version,
				VersionPrerelease: VersionPrerelease,
			}, nil
		},
	}
}

// makeShutdownCh creates an interrupt listener and returns a channel.
// A message will be sent on the channel for every interrupt received.
func makeShutdownCh() <-chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()

	return resultCh
}
