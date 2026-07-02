package main

import (
	"flag"
)

type CliArgs struct {
	Labels   []string
	Up       bool
	Down     bool
	List     bool
	Topology string
}

func GetValidListArgs() []string {
	return []string{"defined", "built", "running", ""}
}

func SetupFlags() *CliArgs {
	cfg := &CliArgs{}

	flag.BoolVar(&cfg.Up, "up", false, "Brings up the specified topology")
	flag.BoolVar(&cfg.Down, "down", false, "Takes down the specified topology")
	flag.StringVar(&cfg.Topology, "topology", "", "Topology name under topologies/<name>; defaults to topologies/default when omitted")
	flag.BoolVar(&cfg.List, "list", false, "Lists all directories available in the topologies directory")

	return cfg
}
