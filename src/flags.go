package main

import (
	"flag"
	"strings"
)

type CliArgs struct {
	Labels   []string
	Start    *bool
	Stop     *bool
	List     *bool
	Topology *string
}

func GetValidListArgs() []string {
	return []string{"defined", "built", "running", ""}
}

func SetupFlags() *CliArgs {
	cfg := &CliArgs{
		Labels: []string{},
	}

	flag.Func("labels", "Scopes commands to the specified labels", func(arg string) error {
		for _, s := range strings.Split(arg, " ") {

			cfg.Labels = append(cfg.Labels, s)
		}

		return nil
	})

	// For later
	// flag.Func("list", fmt.Sprintf("List defined routers and services: [%v]", strings.Join(GetValidListArgs(), ", ")), func(arg string) error {

	// 	if slices.Contains(GetValidListArgs(), arg) {
	// 		cfg.List = &arg
	// 		return nil
	// 	}

	// 	return fmt.Errorf("invalid argument %q: must be one of %v", arg, GetValidListArgs())
	// })

	cfg.Start = flag.Bool("start", false, "Starts the default collection of routers and services")
	cfg.Stop = flag.Bool("stop", false, "Stops all running services")
	cfg.List = flag.Bool("list", false, "Lists all services in the compose file")
	cfg.Topology = flag.String("topology", "", "Topology name under topologies/<name>; defaults to the root compose files when omitted")

	return cfg
}
