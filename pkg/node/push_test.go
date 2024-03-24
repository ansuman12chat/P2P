package node

import (
	"testing"

	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/require"
)

func mockNode(t *testing.T) *Node {
	net := mocknet.New()
	h, err := net.GenPeer()
	require.NoError(t, err)

	return &Node{Host: h}
}
