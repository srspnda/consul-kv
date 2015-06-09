package command

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/consul/api"
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
	var addr, dc, path string

	flagSet := flag.NewFlagSet("get", flag.ContinueOnError)

	flagSet.Usage = func() { g.Ui.Output(g.Help()) }
	flagSet.StringVar(&addr, "http-addr", "127.0.0.1:8500", "http addr of the agent")
	flagSet.StringVar(&dc, "datacenter", "", "datacenter of the agent")
	flagSet.StringVar(&path, "path", "", "path of the kv pair")

	if err := flagSet.Parse(args); err != nil {
		return 0
	}

	if path == "" {
		g.Ui.Error(fmt.Sprintf("Must specify a path to GET"))
		g.Ui.Error("")
		g.Ui.Error(g.Help())
		return 1
	}

	config := &api.Config{
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

	client, err := api.NewClient(config)
	if err != nil {
		g.Ui.Error(fmt.Sprintf("Error connecting to Consul agent: %s", err))
		return 1
	}
	kv := client.KV()

	pair, _, err := kv.Get(path, nil)
	if err != nil {
		g.Ui.Error(fmt.Sprintf("Error retrieving path %s from Consul agent: %s", err))
		return 1
	}

	fmt.Printf("KV: %v", pair)

	return 0
}

// Help returns a string that is the usage for the GetCommand.
func (g *GetCommand) Help() string {
	return strings.TrimSpace(`
Usage: consul-kv get [options]

  Get a KV pair from a Consul agent at the specified path.

Options:

  -http-addr="127.0.0.1:8500"  HTTP address of the Consul agent.
  -datacenter=""               Datacenter of the Consul agent.
  -path=""                     KV pair path to GET from the Consul agent.
`)
}

// Synopsis returns a tring that is the basic usage for the GetCommand.
func (g *GetCommand) Synopsis() string {
	return "Get a KV pair from a local Consul agent at the specified path."
}
