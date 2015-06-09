package command

import (
	"bytes"
	"fmt"

	"github.com/mitchellh/cli"
)

// VersionCommand is a Command implementation that prints `consul-kv` version.
type VersionCommand struct {
	Version string
	Ui      cli.Ui
}

// Run is a function mapped to the VersionCommand implementation.
// This is invoked upon calling `consul-kv [-v|--version|version]`.
func (v *VersionCommand) Run(_ []string) int {
	var version bytes.Buffer

	fmt.Fprintf(&version, "consul-kv v%s", v.Version)
	v.Ui.Output(version.String())

	return 0
}

// Help returns a string that is the usage for the VersionCommand.
func (v *VersionCommand) Help() string {
	return "Prints the consul-kv version"
}

// Synopsis returns a string of the VersionCommand implementation.
func (v *VersionCommand) Synopsis() string {
	return "Prints the consul-kv version"
}
