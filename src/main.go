package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

const rootComposeFile = "topologies/default/docker-compose.yml"
const defaultProfile = "all-isps"
const projectName = "transitlab"

const rootTopologyDir = "../topologies"
const composeFile = "docker-compose.yaml"

const instanceLabelKey = "TransitLabInst"

func main() {

	cfg := SetupFlags()
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	if err := Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func Run(cfg *CliArgs) error {

	if cfg.List {
		o, e := handleList()

		if e != nil {
			fmt.Printf("Error listing topologies: %v", e)
			os.Exit(1)
		}

		fmt.Printf("%s\n", o)
	}

	if cfg.Up {
		o, e := handleUp(cfg.Topology)

		if e != nil {
			fmt.Printf("Error starting topology: %v", e)
			os.Exit(1)
		}

		fmt.Printf("%s\n", o)
	}

	if cfg.Up {
		handleDown("yes")
	}

	return nil
}

func getDockerVersion() {
	fmt.Println("Getting Docker version!")

}

func handleList() (output string, err error) {

	entries, err := os.ReadDir(rootTopologyDir)

	if err != nil {
		return "", fmt.Errorf("Error reading topology directory: %v", err)
	}

	var builder strings.Builder

	for _, entry := range entries {
		if entry.IsDir() {
			builder.WriteString(fmt.Sprintf("➤ %s\n", entry.Name()))
		}
	}

	return builder.String(), nil
}

func handleUp(topologyDir string) (output string, err error) {

	ctx := context.Background()

	dockerCli, err := command.NewDockerCli()

	if err != nil {
		return "", fmt.Errorf("Failed to create Docker CLI: %v", err)
	}

	err = dockerCli.Initialize(&flags.ClientOptions{})

	if err != nil {
		return "", fmt.Errorf("Failed to initialize Docker CLI: %v", err)
	}

	service, err := compose.NewComposeService(dockerCli)

	if err != nil {
		return "", fmt.Errorf("Failed to create a new Docker Compose service instance: %v", err)
	}

	project, err := service.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: []string{fmt.Sprintf("%s/%s/%s", rootTopologyDir, topologyDir, composeFile)},
	})

	if err != nil {
		return "", fmt.Errorf("Failed to load project: %v", err)
	}

	err = service.Up(ctx, project, api.UpOptions{
		Create: api.CreateOptions{},
		Start:  api.StartOptions{},
	})

	if err != nil {
		return "", fmt.Errorf("Failed to start topology: %v", err)
	}

	return fmt.Sprintf("Topology %s started successfully", topologyDir), nil

}

func handleDown(composeFile string) {
	fmt.Println("TBD: Use docker compose directly for now.")
}
