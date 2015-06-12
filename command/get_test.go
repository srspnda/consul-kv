package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/consul/testutil"
	"github.com/mitchellh/cli"
)

func TestGetCommand_implements(t *testing.T) {
	var _ cli.Command = &GetCommand{}
}

func TestGetCommandRun(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Stop()

	srv.SetKV("foo", []byte("bar"))

	ui := new(cli.MockUi)
	c := &GetCommand{Ui: ui}

	args := []string{"-http-addr=" + srv.HTTPAddr, "-key=foo"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad %d. %s", code, ui.ErrorWriter.String())
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "value=bar") {
		t.Fatalf("bad %s", output)
	}
}
