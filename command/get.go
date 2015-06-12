package command

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

// GetCommand is a Command implementation that gets the KV pair from the
// peer set of a Consul agent.
type GetCommand struct {
	Ui cli.Ui
}

// Run is the function mapped to the GetCommand implmentation.
// This function is called upon execution of `consul-kv get ...`.
func (g *GetCommand) Run(args []string) int {
	var addr, dc, key string

	flagSet := flag.NewFlagSet("get", flag.ContinueOnError)

	flagSet.Usage = func() { g.Ui.Output(g.Help()) }
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
		g.Ui.Error(fmt.Sprintf("-key required"))
		g.Ui.Error("")
		g.Ui.Error(g.Help())
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

	c, err := consulapi.NewClient(config)
	if err != nil {
		g.Ui.Error(fmt.Sprintf("%s", err))
		return 1
	}
	kv := c.KV()

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		g.Ui.Error(fmt.Sprintf("%s", err))
		return 1
	}

	g.Ui.Output(fmt.Sprintf("GET: key=%s value=%s\n", key, pair.Value))

	return 0
}

// Help returns a string that is the usage for the GetCommand.
func (g *GetCommand) Help() string {
	return strings.TrimSpace(`
Usage: consul-kv get [options]

  Get a KV pair from a Consul agent at the specified key path.

Options:

  -http-addr="127.0.0.1:8500"  HTTP address of the Consul agent.
  -datacenter=""               Datacenter of the Consul agent.
  -key=""                      Key of KV pair to GET from the Consul agent.
`)
}

// Synopsis returns a tring that is the basic usage for the GetCommand.
func (g *GetCommand) Synopsis() string {
	return "Get a KV pair from a Consul agent at the specified key path."
}
