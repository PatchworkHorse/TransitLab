package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/compose-spec/compose-go/v2/cli"
)

const composeFile = "docker-compose.yml"
const defaultProfile = "all-isps"
const projectName = "internetemulator"

func main() {

	cfg := SetupFlags()
	flag.Parse()

	if err := Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func Run(cfg *CliArgs) error {

	activeStateArgs := 0

	for _, f := range []*bool{cfg.Start, cfg.Stop} {
		if *f {
			activeStateArgs++
		}
	}

	if activeStateArgs > 1 {
		return fmt.Errorf("Only one state-modifying command may specified")
	}

	if *cfg.List {
		handleList()
	}

	if *cfg.Start {
		handleStart(defaultProfile)
	}

	if *cfg.Stop {
		handleStop()
	}

	return nil
}

func handleList() {

	ctx := context.Background()

	options, err := cli.NewProjectOptions(
		[]string{composeFile},
		cli.WithWorkingDirectory("../"),
		cli.WithOsEnv,
		cli.WithDotEnv,
		cli.WithName("test"),
		cli.WithProfiles([]string{"*"}),
	)
	if err != nil {
		log.Fatal(err)
	}

	project, err := options.LoadProject(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, service := range project.Services {
		if _, ok := service.Labels["horse.patchwork.netemu.router"]; !ok {
			continue
		}

		fmt.Printf("Service: %s, Profiles: %v\n", service.Name, service.Profiles)

		if _, ok := service.Labels["horse.patchwork.netemu.template"]; ok {
			continue
		}

		fmt.Println(service.Name)
		fmt.Printf("  Interfaces:\n")
		for netName, net := range service.Networks {
			fmt.Printf("    %s -> %s\n", net.InterfaceName, netName)
		}

	}

}

func handleStart(profile string) {
	cmd := exec.Command("docker", "compose", "-f", composeFile, "--profile", profile, "up", "--build", "--quiet-build", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("docker compose up failed: %v", err)
	}
	log.Println("Services started successfully!")
}

func handleStop() {
	cmd := exec.Command("docker", "compose", "-p", projectName, "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("docker compose down failed: %v", err)
	}

	log.Println("Services stopped successfully!")
}
