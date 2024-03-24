package node

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/ansuman12chat/p2p/pkg/config"
	p2p "github.com/ansuman12chat/p2p/pkg/pb"
)

// authenticateMessages is used in tests to skip
// message authentication due to bogus keys.
var authenticateMessages = true

// Node encapsulates the logic for sending and receiving messages.
type Node struct {
	host.Host
	*MDNSProtocol
	*PushProtocol
	*TransferProtocol
}

// Init creates a new, fully initialized node with the given options.
func Init(ctx context.Context, opts ...libp2p.Option) (*Node, error) {

	conf, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	if !conf.Identity.IsInitialized() {
		err = conf.Identity.GenerateKeyPair()
		if err != nil {
			return nil, err
		}
	}

	key, err := conf.Identity.PrivateKey()
	if err != nil {
		return nil, err
	}
	opts = append(opts, libp2p.Identity(key))
	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	node := &Node{Host: h}
	node.MDNSProtocol = NewMDNSProtocol(node)
	node.PushProtocol = NewPushProtocol(node)
	node.TransferProtocol = NewTransferProtocol(node)

	return node, nil
}

// Send prepares the message msg to be sent over the network stream s.
// Send closes the stream for writing but leaves it open for reading.
func (n *Node) Send(s network.Stream, msg p2p.HeaderMessage) error {
	defer s.CloseWrite()

	// Get own public key.
	pubKey := n.Host.Peerstore().PubKey(n.Host.ID())
	pubKeyBytes, err := crypto.MarshalPublicKey(pubKey)
	if err != nil {
		return err
	}

	hdr := &p2p.Header{
		RequestId:  uuid.New().String(),
		NodeId:     n.Host.ID().String(),
		NodePubKey: pubKeyBytes,
		Timestamp:  appTime.Now().Unix(),
	}
	msg.SetHeader(hdr)

	// Transform msg to binary to calculate the signature.
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	// Sign the data and attach the signature.
	key := n.Host.Peerstore().PrivKey(n.Host.ID())
	signature, err := key.Sign(data)
	if err != nil {
		return err
	}
	hdr.Signature = signature
	msg.SetHeader(hdr) // Maybe unnecessary

	// Transform msg + signature to binary.
	data, err = proto.Marshal(msg)
	if err != nil {
		return err
	}

	// Transmit the data.
	_, err = s.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// authenticateMessage verifies the authenticity of the message payload.
// It takes the given signature and verifies it against the given public
// key. It returns true if the signature is valid, false otherwise.
func (n *Node) authenticateMessage(msg p2p.HeaderMessage) (bool, error) {

	// Short circuit in test runs with invalid keys.
	if !authenticateMessages {
		return true, nil
	}

	// store a temp ref to signature and remove it from message msg
	// sign is a string to allow easy reset to zero-value (empty string)
	signature := msg.GetHeader().Signature
	msg.GetHeader().Signature = nil

	// marshall msg without the signature to protobufs3 binary format
	bin, err := proto.Marshal(msg)
	if err != nil {
		return false, err
	}

	// restore sig in message msg (for possible future use)
	msg.GetHeader().Signature = signature

	// restore peer id binary format from base58 encoded node id msg
	peerID, err := peer.Decode(msg.GetHeader().NodeId)
	if err != nil {
		return false, err
	}

	key, err := crypto.UnmarshalPublicKey(msg.GetHeader().NodePubKey)
	if err != nil {
		return false, fmt.Errorf("failed to extract key from message key msg")
	}

	// extract node id from the provided public key
	idFromKey, err := peer.IDFromPublicKey(key)

	if err != nil {
		return false, fmt.Errorf("failed to extract peer id from public key")
	}

	// verify that message author node id matches the provided node public key
	if idFromKey != peerID {
		return false, fmt.Errorf("node id and provided public key mismatch")
	}

	return key.Verify(bin, signature)
}

// Read drains the given stream and parses the content. It unmarshalls
// it into the protobuf object. It also verifies the authenticity of the message.
// Read closes the stream for reading but leaves it open for writing.
func (n *Node) Read(s network.Stream, data p2p.HeaderMessage) error {
	defer s.CloseRead()
	buf, err := io.ReadAll(s)
	if err != nil {
		if err2 := s.Reset(); err2 != nil {
			err = errors.Wrap(err, err2.Error())
		}
		return err
	}

	if err = proto.Unmarshal(buf, data); err != nil {
		return err
	}

	valid, err := n.authenticateMessage(data)
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("failed to authenticate message")
	}

	return nil
}

// WaitForEOF waits for the EOF signal on the given stream.
func (n *Node) WaitForEOF(s network.Stream) error {
	// 10 sec timeout
	timeout := time.After(10 * time.Second)
	done := make(chan error)
	go func() {
		buf := make([]byte, 1)
		n, err := s.Read(buf)
		if err == io.EOF && n == 0 {
			err = nil
		} else if n != 0 {
			err = fmt.Errorf("stream returned data unexpectedly")
		}
		done <- err
		close(done)
	}()
	select {
	case <-timeout:
		return fmt.Errorf("timeout")
	case err := <-done:
		return err
	}
	// Select on the done channel and a timeout channel.
}
