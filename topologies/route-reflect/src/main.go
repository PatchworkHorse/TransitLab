package main

import (
	"context"
	"fmt"
	"log"

	api "github.com/osrg/gobgp/v4/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Patchwork notes:
// Some experiments with the GoBGP gRPC API. Brainstorm more ideas for what we can do here.

func main() {
	fmt.Println("Connecting to GoBGP at localhost:50051...")
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connected! Creating client...")

	client := api.NewGoBgpServiceClient(conn)
	ctx := context.Background()

	fmt.Println("Starting watcher...")
	stream, err := watchNeighbors(ctx, client)
	if err != nil {
		log.Fatalf("failed to start watcher: %v", err)
	}
	fmt.Println("Watcher started, waiting for events...")

	for event := range stream {
		if event.Peer != nil && event.Peer.State != nil {
			fmt.Printf("neighbor %s state=%s\n",
				event.Peer.State.NeighborAddress,
				event.Peer.State.SessionState,
			)
		}
	}
}

func watchNeighbors(ctx context.Context, client api.GoBgpServiceClient) (<-chan *api.WatchEventResponse_PeerEvent, error) {
	stream, err := client.WatchEvent(ctx, &api.WatchEventRequest{
		Peer: &api.WatchEventRequest_Peer{},
	})
	if err != nil {
		return nil, err
	}

	ch := make(chan *api.WatchEventResponse_PeerEvent, 16)
	go func() {
		defer close(ch)
		for {
			resp, err := stream.Recv()
			if err != nil {
				return
			}
			if peer := resp.GetPeer(); peer != nil {
				ch <- peer
			}
		}
	}()

	return ch, nil
}
