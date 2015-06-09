package main

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/srspnda/consul-kv/command"
)

// Commands is the mapping of all the available consul-kv commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	Commands = map[string]cli.CommandFactory{
		"get": func() (cli.Command, error) {
			return &command.GetCommand{
				Ui: ui,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Version: Version,
				Ui:      ui,
			}, nil
		},
	}
}
