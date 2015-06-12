package command

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

// PutCommand is a Command implementation that puts the KV pair on the
// peer set of a Consul agent.
type PutCommand struct {
	Ui cli.Ui
}

// Run is the function mapped to the PutCommand implmentation.
// This function is called upon execution of `consul-kv put ...`.
func (p *PutCommand) Run(args []string) int {
	var addr, dc, key, value string

	flagSet := flag.NewFlagSet("put", flag.ContinueOnError)
	flagSet.Usage = func() { p.Ui.Output(p.Help()) }
	flagSet.StringVar(
		&addr,
		"http-addr",
		"127.0.0.1:8500",
		"http addr of the agent",
	)
	flagSet.StringVar(&dc, "datacenter", "", "datacenter of the agent")
	flagSet.StringVar(&key, "key", "", "key of the kv pair")
	flagSet.StringVar(&value, "value", "", "value of the kv pair")

	if err := flagSet.Parse(args); err != nil {
		return 0
	}

	if key == "" || value == "" {
		p.Ui.Error(fmt.Sprintf("-key and -value required"))
		p.Ui.Error("")
		p.Ui.Error(p.Help())
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
		p.Ui.Error(fmt.Sprintf("%s", err))
		return 1
	}
	kv := client.KV()

	pair := &consulapi.KVPair{Key: key, Value: []byte(value)}
	if _, err := kv.Put(pair, nil); err != nil {
		p.Ui.Error(fmt.Sprintf("%s", err))
		return 1
	}

	p.Ui.Output(fmt.Sprintf("PUT: key=%s value=%s\n", key, value))

	return 0
}

// Help returns a string that is the usage for the PutCommand.
func (p *PutCommand) Help() string {
	return strings.TrimSpace(`
Usage: consul-kv put [options]

  Put a KV pair from a Consul agent at the specified key path and value.

Options:

  -http-addr="127.0.0.1:8500"  HTTP address of the Consul agent.
  -datacenter=""               Datacenter of the Consul agent.
  -key=""                      Key path to PUT on the Consul agent.
  -value=""                    Value of key path to PUT on the Consul agent.
`)
}

// Synopsis returns a tring that is the basic usage for the PutCommand.
func (p *PutCommand) Synopsis() string {
	return "Put a KV pair on a Consul agent at the specified key path."
}
