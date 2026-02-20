package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
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

	if *cfg.List {
		handleList()
	}

	if *cfg.BgpSum != "" {
		handleBgpSummary(*cfg.BgpSum)
	}

	if *cfg.Start {
		handleStart(defaultProfile)
	}

	if *cfg.Stop {
		handleStop()
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

// Todo: Remove hardcoded labels
func handleList() {

	for _, service := range getProject().Services {
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

func handleBgpSummary(routerName string) {

	routerCmd := "sho ip bgp summary json"
	args := []string{"exec", routerName, "vtysh", "-c", routerCmd}
	cmd := exec.Command("docker", args...)

	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute BGP summary commands: %v", err)
	}

	var summary BGPSummary

	json.Unmarshal(output, &summary)

	fmt.Printf("🏘️  BGP Peering summary for %s:\n", routerName)
	fmt.Printf("%d peers found\n", summary.IPv4Unicast.PeerCount)

	for peerName, peer := range summary.IPv4Unicast.Peers {
		fmt.Printf("  Peer: %s | %s | %d\n", peerName, peer.Hostname, peer.RemoteAS)
		fmt.Printf("    Local AS: %d\n", peer.LocalAS)
		fmt.Printf("    Peer Uptime: %s\n", peer.PeerUptime)

		// Todo: Just switch for the status bubble, move to helper
		switch peer.State {
		case "Established":
			fmt.Printf("    State: %s 🟢\n", peer.State)
		case "Active":
			fmt.Printf("    State: %s 🔵\n", peer.State)
		case "Idle":
			fmt.Printf("    State: %s 🔴", peer.State)
		}
	}
}

func resolveServiceByName(serviceName string) (types.ServiceConfig, bool) {

	if service, ok := getProject().Services[serviceName]; ok {
		return service, true
	}
	return types.ServiceConfig{}, false
}

func getProject() *types.Project {
	ctx := context.Background()

	options, err := cli.NewProjectOptions(
		[]string{composeFile},
		cli.WithWorkingDirectory("."),
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

	return project
}
