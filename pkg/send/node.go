package send

import (
	"context"
	"os"
	"path"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"

	"github.com/ansuman12chat/p2p/internal/log"
	"github.com/ansuman12chat/p2p/pkg/node"
	"github.com/ansuman12chat/p2p/pkg/progress"
)

type Node struct {
	*node.Node
}

func InitNode(ctx context.Context) (*Node, error) {

	n, err := node.Init(ctx)
	if err != nil {
		return nil, err
	}

	return &Node{n}, nil
}

func (n *Node) Close() error {
	err := n.Host.Close()
	if err != nil {
		log.Infoln(err)
	}

	err = n.StopMdnsService()
	if err != nil {
		log.Infoln(err)
	}

	return nil
}

func (n *Node) Transfer(ctx context.Context, pi peer.AddrInfo, filepath string) (bool, error) {
	// Connect to peer
	err := n.Connect(ctx, pi)
	if err != nil {
		return false, err
	}
	// Get content ID
	c, err := calcContentID(filepath)
	if err != nil {
		return false, err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Get file info
	fstat, err := f.Stat()
	if err != nil {
		return false, err
	}
	log.Infof("Asking for confirmation... ")

	accepted, err := n.SendPushRequest(ctx, pi.ID, path.Base(f.Name()), fstat.Size(), c)
	if err != nil {
		return false, err
	}

	if !accepted {
		log.Infoln("Rejected!")
		return accepted, nil
	}
	log.Infoln("Accepted!")

	pr := progress.NewReader(f)

	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	go node.IndicateProgress(ctx, pr, path.Base(f.Name()), fstat.Size(), &wg)
	defer func() { cancel(); wg.Wait() }()

	if _, err = n.Node.Transfer(ctx, pi.ID, pr); err != nil {
		return accepted, errors.Wrap(err, "could not transfer file to peer")
	}

	log.Infoln("Successfully sent file!")
	return accepted, nil
}
