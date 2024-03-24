package node

import (
	"bytes"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/assert"

	"github.com/ansuman12chat/p2p/internal/log"
	"github.com/ansuman12chat/p2p/internal/mock"
)

func TestTransferProtocol_onTransfer_unexpected(t *testing.T) {
	n := mockNode(t)
	p := NewTransferProtocol(n)

	ctrl := gomock.NewController(t)
	streamer := mock.NewMockStreamer(ctrl)
	conner := mock.NewMockConner(ctrl)
	conner.EXPECT().RemotePeer().Return(peer.ID("peer-id"))

	streamer.EXPECT().Conn().Return(conner)
	streamer.EXPECT().Reset().Return(nil)

	buffer := new(bytes.Buffer)
	log.Out = buffer
	p.onTransfer(streamer)

	logOut := strings.TrimSpace(buffer.String())

	assert.Equal(t, logOut, "Received data transfer attempt from unexpected peer")
}
