package command

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

// SetCommand is a Command implementation that sets the KV pair on the
// peer set of a Consul agent.
type SetCommand struct {
	Ui cli.Ui
}

// Run is the function mapped to the SetCommand implmentation.
// This function is called upon execution of `consul-kv set ...`.
func (s *SetCommand) Run(args []string) int {
	var addr, dc, key, value string

	flagSet := flag.NewFlagSet("set", flag.ContinueOnError)
	flagSet.Usage = func() { s.Ui.Output(s.Help()) }
	flagSet.StringVar(&addr, "http-addr", "127.0.0.1:8500", "http addr of the agent")
	flagSet.StringVar(&dc, "datacenter", "", "datacenter of the agent")
	flagSet.StringVar(&key, "key", "", "key of the kv pair")
	flagSet.StringVar(&value, "value", "", "value of the kv pair")

	if err := flagSet.Parse(args); err != nil {
		return 0
	}

	if key == "" || value == "" {
		s.Ui.Error(fmt.Sprintf("-key and -value required"))
		s.Ui.Error("")
		s.Ui.Error(s.Help())
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
		s.Ui.Error(fmt.Sprintf("error on client addr=%s dc=%s: %s\n", addr, dc, err))
		return 1
	}
	kv := client.KV()

	p := &consulapi.KVPair{Key: key, Value: []byte(value)}
	if _, err := kv.Put(p, nil); err != nil {
		s.Ui.Error(fmt.Sprintf("error on SET key=%s value=%s: %s\n", key, value, err))
		return 1
	}

	fmt.Printf("SET: key=%s value=%s\n", key, value)

	return 0
}

// Help returns a string that is the usage for the SetCommand.
func (g *SetCommand) Help() string {
	return strings.TrimSpace(`
Usage: consul-kv set [options]

  Set a KV pair from a Consul agent at the specified key path and value.

Options:

  -http-addr="127.0.0.1:8500"  HTTP address of the Consul agent.
  -datacenter=""               Datacenter of the Consul agent.
  -key=""                      Key path to SET on the Consul agent.
  -value=""                    Value of key path to SET on the Consul agent.
`)
}

// Synopsis returns a tring that is the basic usage for the SetCommand.
func (g *SetCommand) Synopsis() string {
	return "Set a KV pair on a Consul agent at the specified key path."
}
