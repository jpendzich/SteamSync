package network

import (
	"os"
	"time"

	"github.com/schollz/peerdiscovery"
)

type Peer struct {
	IpAdress string
	Hostname string
}

func GetAllPeers() ([]Peer, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	discovered, err := peerdiscovery.Discover(peerdiscovery.Settings{
		Payload:   []byte(hostname),
		TimeLimit: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	peers := make([]peer, len(discovered))
	for i := 0; i < len(peers); i += 1 {
		peers[i].IpAdress = discovered[i].Address
		peers[i].Hostname = string(discovered[i].Payload)
	}
	return peers, nil
}
