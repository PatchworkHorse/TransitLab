package main

type BGPSummary struct {
	IPv4Unicast struct {
		PeerCount int                `json:"peerCount"`
		Peers     map[string]BGPPeer `json:"peers"`
	} `json:"ipv4Unicast"`
}

type BGPPeer struct {
	Hostname   string `json:"hostname,omitempty"`
	RemoteAS   int    `json:"remoteAs"`
	LocalAS    int    `json:"localAs"`
	Version    int    `json:"version"`
	PeerUptime string `json:"peerUptime"`
	State      string `json:"state"`
}
