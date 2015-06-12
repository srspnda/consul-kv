package command

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

// LsCommand is a Command implementation that lists the keys on the
// peer set of a Consul agent at a specified key prefix.
type LsCommand struct {
	Ui cli.Ui
}

// Run is the function mapped to the LsCommand implmentation.
// This function is called upon execution of `consul-kv ls ...`.
func (l *LsCommand) Run(args []string) int {
	var addr, dc, key string
	flagSet := flag.NewFlagSet("ls", flag.ContinueOnError)
	flagSet.Usage = func() { l.Ui.Output(l.Help()) }
	flagSet.StringVar(
		&addr,
		"http-addr",
		"127.0.0.1:8500",
		"http addr of the agent",
	)
	flagSet.StringVar(&dc, "datacenter", "", "datacenter of the agent")
	flagSet.StringVar(&key, "key", "", "key of the kv pair")
	if err := flagSet.Parse(args); err != nil {
		return 0
	}
	if key == "" {
		l.Ui.Error(fmt.Sprintf("-key required"))
		l.Ui.Error("")
		l.Ui.Error(l.Help())
		return 1
	}
	config := &consulapi.Config{
		Address:    "127.0.0.1:8500",
		Scheme:     "http",
		HttpClient: http.DefaultClient,
	}
	if addr != "" {
		config.Address = addr
	}
	if dc != "" {
		config.Datacenter = dc
	}
	client, err := consulapi.NewClient(config)
	if err != nil {
		l.Ui.Error(fmt.Sprintf("%s", err))
		return 1
	}
	kv := client.KV()

	entries, _, err := kv.List(key, nil)
	if err != nil {
		l.Ui.Error(fmt.Sprintf("%s", err))
		return 1
	}

	var keys []string
	for _, e := range entries {
		keys = append(keys, e.Key)
	}

	l.Ui.Output(fmt.Sprintf("LS: keys={%s}", strings.Join(keys, ", ")))

	return 0
}

// Help returns a string that is the usage for the LsCommand.
func (l *LsCommand) Help() string {
	return strings.TrimSpace(`
Usage: consul-kv put [options]

  List all keys on a Consul agent under the specified key path.

Options:

  -http-addr="127.0.0.1:8500"  HTTP address of the Consul agent.
  -datacenter=""               Datacenter of the Consul agent.
  -key=""                      Key path to list all subkeys on the Consul agent.
`)
}

// Synopsis returns a tring that is the basic usage for the LsCommand.
func (l *LsCommand) Synopsis() string {
	return "List all keys on a Consul agent under the specified key path."
}
