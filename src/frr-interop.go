package main

type BGPStatus struct {
	IPv4Unicast struct {
		PeerCount int `json:"peerCount"`
		Peers     map[string]struct {
			Hostname   string `json:"hostname,omitempty"`
			RemoteAS   int    `json:"remoteAs"`
			LocalAS    int    `json:"localAs"`
			Version    int    `json:"version"`
			PeerUptime string `json:"peerUptime"`
			State      string `json:"state"`
		} `json:"peers"`
	} `json:"ipv4Unicast"`
}
