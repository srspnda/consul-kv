package command

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/consul/testutil"
	"github.com/mitchellh/cli"
)

func TestLsCommand_implements(t *testing.T) {
	var _ cli.Command = &LsCommand{}
}

func TestLsCommandRun(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Stop()

	srv.PopulateKV(map[string][]byte{
		"foo/bar":  []byte("baz"),
		"foo/man":  []byte("shu"),
		"foo/test": []byte("123"),
	})

	ui := new(cli.MockUi)
	c := &LsCommand{Ui: ui}

	args := []string{"-http-addr=" + srv.HTTPAddr, "-key=foo/"}
	code := c.Run(args)
	if code != 0 {
		t.Fatalf("bad %d. %s", code, ui.ErrorWriter.String())
	}

	output := ui.OutputWriter.String()
	fmt.Println(output)
	if !strings.Contains(output, "foo/bar, foo/man, foo/test") {
		t.Fatalf(
			"bad ls. keys={foo/bar, foo/man, foo/test} got keys={%s}",
			output,
		)
	}
}
