package node

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"

	"github.com/ansuman12chat/p2p/internal/app"
	"github.com/ansuman12chat/p2p/internal/log"
	commons "github.com/ansuman12chat/p2p/pkg/commons"
)

// The time a discovered peer will stay in the `Peers` map.
// If it doesn't get re-discovered within this time it gets
// removed from the list.
var gcDuration = 5 * time.Second

// Variable assignments for mocking purposes.
var (
	appTime app.Timer = app.Time{}
)

// MDNSProtocol encapsulates the logic for discovering peers
// in via multicast DNS in the local network.
type MDNSProtocol struct {
	node         *Node
	mdnsServ     mdns.Service
	MdnsInterval time.Duration
	Peers        *sync.Map
}

// NewMdnsProtocol creates a new MDNSProtocol struct with
// sane defaults.
func NewMdnsProtocol(node *Node) *MDNSProtocol {
	m := &MDNSProtocol{
		node:         node,
		MdnsInterval: time.Second,
		Peers:        &sync.Map{},
	}
	return m
}

// StartMdnsService starts the mDNS service and registers the
// MDNSProtocol to be notified for every newly discovered peer.
func (m *MDNSProtocol) StartMdnsService(ctx context.Context) error {
	if m.mdnsServ != nil {
		return nil
	}

	mdns := mdns.NewMdnsService(m.node.Host, commons.ServiceTag, m)
	m.mdnsServ = mdns
	err := m.mdnsServ.Start()
	if err != nil {
		log.Infof("Starting the mDNS service failed: %s", err)
		return err
	}
	return nil
}

// StopMdnsService stops the mDNS service and clears the list
// of peers.
func (m *MDNSProtocol) StopMdnsService() error {
	if m.mdnsServ == nil {
		return nil
	}

	err := m.mdnsServ.Close()
	if err != nil {
		return err
	}

	m.Peers.Range(func(key, value interface{}) bool {
		value.(PeerInfo).timer.Stop()
		return true
	})

	m.mdnsServ = nil
	// Clearning the list of peers
	m.Peers = &sync.Map{}

	return nil
}

type PeerInfo struct {
	pi    peer.AddrInfo
	timer *time.Timer
}

// HandlePeerFound stores every newly found peer in a map.
// Every map entry gets a timer assigned that removes the
// entry after a garbage collection timeout if the peer
// is not seen again in the meantime. If we see the peer
// again we reset the time to start again from that point
// in time.
func (m *MDNSProtocol) HandlePeerFound(pi peer.AddrInfo) {
	savedPeer, ok := m.Peers.Load(pi.ID)
	if ok {
		savedPeer.(PeerInfo).timer.Reset(gcDuration)
	} else {
		// If the peer is not in the list, add it with a timer
		t := appTime.AfterFunc(gcDuration, func() {
			m.Peers.Delete(pi.ID)
		})
		m.Peers.Store(pi.ID, PeerInfo{pi, t})
	}
}

// PeersList returns a sorted list of address information
// structs. Sorting order is based on the peer ID.
func (m *MDNSProtocol) PeersList() []peer.AddrInfo {
	peers := []peer.AddrInfo{}
	m.Peers.Range(func(key, value interface{}) bool {
		peers = append(peers, value.(PeerInfo).pi)
		return true
	})

	sort.Slice(peers, func(i, j int) bool {
		return peers[i].ID < peers[j].ID
	})

	return peers
}

// PrintPeers dumps the given list of peers to the screen
// to be selected by the user via its index.
func (m *MDNSProtocol) PrintPeers(peers []peer.AddrInfo) {
	for i, p := range peers {
		fmt.Fprintf(os.Stdout, "[%d] %s\n", i, p.ID)
	}
	fmt.Fprintln(os.Stdout)
}
