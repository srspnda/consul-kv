package command

import (
	"bytes"
	"testing"

	"github.com/hashicorp/consul/testutil"
	"github.com/mitchellh/cli"
)

func TestPutCommand_implements(t *testing.T) {
	var _ cli.Command = &GetCommand{}
}

func TestPutCommandRun(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Stop()

	ui := new(cli.MockUi)
	c := &PutCommand{Ui: ui}

	args := []string{
		"-http-addr=" + srv.HTTPAddr,
		"-key=foo",
		"-value=bar",
	}

	code := c.Run(args)
	if code != 0 {
		t.Fatalf("bad %d. %s", code, ui.ErrorWriter.String())
	}

	value := srv.GetKV("foo")
	if !bytes.Equal(value, []byte("bar")) {
		t.Fatalf("bad put. key=foo value=bar, not value=%s", value)
	}
}
