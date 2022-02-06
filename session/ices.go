package session

import "github.com/felixqin/p2p-tunnel/tunnel"

var iceServers tunnel.IceServers

func AddIceServer(server string) {
	iceServers = append(iceServers, server)
}

func IceServers() *tunnel.IceServers {
	return &iceServers
}
