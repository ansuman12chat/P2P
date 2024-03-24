package node

import (
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/assert"

	"github.com/ansuman12chat/p2p/internal/app"
	"github.com/ansuman12chat/p2p/internal/mock"
)

func setup(t *testing.T) *gomock.Controller {
	appTime = app.Time{}
	return gomock.NewController(t)
}

func teardown(t *testing.T, ctrl *gomock.Controller) {
	ctrl.Finish()
	appTime = app.Time{}
}

func TestNewMDNSProtocol_returnsInitializedMDNSProtocol(t *testing.T) {
	n := mockNode(t)
	m := NewMDNSProtocol(n)
	assert.Equal(t, n, m.node)
	assert.Equal(t, time.Second, m.MdnsInterval)
	assert.NotNil(t, m.Peers)
}

func TestMDNSProtocol_HandlePeerFound_storesPeer(t *testing.T) {
	n := mockNode(t)
	p := NewMDNSProtocol(n)

	pi := peer.AddrInfo{ID: peer.ID("peer-id")}
	p.HandlePeerFound(pi)

	_, found := p.Peers.Load(peer.ID("peer-id"))
	assert.True(t, found)
}

func TestMDNSProtocol_HandlePeerFound_deletesPeerAfterGCTime(t *testing.T) {
	n := mockNode(t)
	p := NewMDNSProtocol(n)

	ctrl := setup(t)
	defer teardown(t, ctrl)

	gcDurationTmp := gcDuration
	gcDuration = 0 * time.Millisecond
	peerID := peer.ID("peer-id")

	var wg sync.WaitGroup
	wg.Add(1)

	mTime := mock.NewMockTimer(ctrl)
	mTime.EXPECT().
		AfterFunc(gomock.Eq(gcDuration), gomock.Any()).
		Times(1).
		DoAndReturn(func(d time.Duration, f func()) *time.Timer {
			return time.AfterFunc(d, func() {
				f()
				wg.Done()
			})
		})

	appTime = mTime

	pi := peer.AddrInfo{ID: peerID}
	p.HandlePeerFound(pi)

	wg.Wait()

	_, found := p.Peers.Load(peerID)
	assert.False(t, found)

	gcDuration = gcDurationTmp
}

func TestMDNSProtocol_HandlePeerFound_deletesPeerAfterGCTimeWithIntermediateResets(t *testing.T) {
	n := mockNode(t)
	p := NewMDNSProtocol(n)

	ctrl := setup(t)
	defer teardown(t, ctrl)

	gcDurationTmp := gcDuration
	gcDuration = 100 * time.Millisecond
	peerID := peer.ID("peer-id")

	mTime := mock.NewMockTimer(ctrl)
	mTime.EXPECT().
		AfterFunc(gomock.Eq(gcDuration), gomock.Any()).
		Times(1).
		DoAndReturn(time.AfterFunc)
	appTime = mTime

	pi := peer.AddrInfo{ID: peerID}
	p.HandlePeerFound(pi)
	time.Sleep(50 * time.Millisecond)
	p.HandlePeerFound(pi)
	time.Sleep(75 * time.Millisecond)
	_, found := p.Peers.Load(peerID)
	assert.True(t, found)
	time.Sleep(50 * time.Millisecond)
	_, found = p.Peers.Load(peerID)
	assert.False(t, found)

	gcDuration = gcDurationTmp
}

func TestMDNSProtocol_PeerList_returnsListOfPeers(t *testing.T) {
	n := mockNode(t)
	p := NewMDNSProtocol(n)

	p1 := peer.AddrInfo{ID: peer.ID("peer-id-1")}
	p2 := peer.AddrInfo{ID: peer.ID("peer-id-2")}
	p3 := peer.AddrInfo{ID: peer.ID("peer-id-3")}

	p.HandlePeerFound(p1)
	p.HandlePeerFound(p3)
	p.HandlePeerFound(p2)

	list := p.PeersList()

	assert.Equal(t, p1, list[0])
	assert.Equal(t, p2, list[1])
	assert.Equal(t, p3, list[2])
}
