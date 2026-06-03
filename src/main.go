package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/compose-spec/compose-go/v2/cli"
)

const rootComposeFile = "docker-compose.yml"
const defaultProfile = "all-isps"
const projectName = "transitlab"

func main() {

	cfg := SetupFlags()
	flag.Parse()

	if err := Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func Run(cfg *CliArgs) error {

	if dockerVersion, error := getDockerVersion(); error == nil {
		fmt.Printf("Found Docker: %s\n", dockerVersion)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", error)
		os.Exit(1)
	}

	activeStateArgs := 0

	for _, f := range []*bool{cfg.Start, cfg.Stop} {
		if *f {
			activeStateArgs++
		}
	}

	if activeStateArgs > 1 {
		return fmt.Errorf("Only one state-modifying command may specified")
	}

	composeFile := resolveComposeFile(*cfg.Topology)
	project := resolveProjectName(*cfg.Topology)

	if *cfg.List {
		handleList(composeFile)
	}

	if *cfg.Start {
		handleStart(composeFile, project, defaultProfile)
	}

	if *cfg.Stop {
		handleStop(composeFile, project)
	}

	return nil
}

func getDockerVersion() (string, error) {
	cmd := exec.Command("docker", "version", "--format", "{{.Client.Version}}")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get docker version: %w", err)
	}
	return string(output), nil
}

func resolveComposeFile(topology string) string {
	topology = strings.TrimSpace(topology)
	if topology == "" {
		return rootComposeFile
	}

	return filepath.Join("topologies", topology, "docker-compose.yml")
}

func resolveProjectName(topology string) string {
	topology = strings.TrimSpace(topology)
	if topology == "" {
		return projectName
	}

	clean := regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(strings.ToLower(topology), "-")
	clean = strings.Trim(clean, "-")
	if clean == "" {
		clean = "default"
	}

	return fmt.Sprintf("%s-%s", projectName, clean)
}

func handleList(composeFile string) {

	ctx := context.Background()
	workingDir := filepath.Dir(composeFile)
	if workingDir == "." || workingDir == "" {
		workingDir = "."
	}

	options, err := cli.NewProjectOptions(
		[]string{composeFile},
		cli.WithWorkingDirectory(workingDir),
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
		if _, ok := service.Labels["horse.patchwork.transitlab.router"]; !ok {
			continue
		}

		fmt.Printf("Service: %s, Profiles: %v\n", service.Name, service.Profiles)

		if _, ok := service.Labels["horse.patchwork.transitlab.template"]; ok {
			continue
		}

		fmt.Println(service.Name)
		fmt.Printf("  Interfaces:\n")
		for netName, net := range service.Networks {
			fmt.Printf("    %s -> %s\n", net.InterfaceName, netName)
		}

	}

}

func handleStart(composeFile string, project string, profile string) {
	cmd := exec.Command("docker", "compose", "-f", composeFile, "-p", project, "--profile", profile, "up", "--build", "--quiet-build", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("docker compose up failed: %v", err)
	}
	log.Println("Services started successfully!")
}

func handleStop(composeFile string, project string) {
	cmd := exec.Command("docker", "compose", "-f", composeFile, "-p", project, "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("docker compose down failed: %v", err)
	}

	log.Println("Services stopped successfully!")
}
