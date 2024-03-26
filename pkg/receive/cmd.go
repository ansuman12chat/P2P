package receive

import (
	"encoding/hex"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/ansuman12chat/p2p/internal/log"
	"github.com/ansuman12chat/p2p/pkg/config"
	p2p "github.com/ansuman12chat/p2p/pkg/pb"
)

var Command = &cli.Command{
	Name:    "receive",
	Usage:   "waits until a peer attempts to connect in your local network to receive your a file",
	Aliases: []string{"r"},
	Action:  Action,
	Flags: []cli.Flag{
		&cli.Int64Flag{
			Name:    "port",
			EnvVars: []string{"P2P_PORT"},
			Aliases: []string{"p"},
			Usage:   "The port at which you are reachable for other peers in the network.",
			Value:   44044,
		},
		&cli.StringFlag{
			Name:    "host",
			EnvVars: []string{"P2P_HOST"},
			Usage:   "The host at which you are reachable for other peers in the network.",
			Value:   "0.0.0.0",
		},
	},
	ArgsUsage:   "[DEST_DIR]",
	UsageText:   ``,
	Description: `The receive subcommand will wait for a peer to connect to your node and receive a file.`,
}

// Action is the function that is called when running p2p receive.
func Action(c *cli.Context) error {
	shutdown := make(chan error)

	ctx, err := config.FillContext(c.Context)
	if err != nil {
		return errors.Wrap(err, "failed loading configuration")
	}

	local, err := InitNode(ctx, c.String("host"), c.Int64("port"), shutdown)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to initialize node"))
	}
	defer local.Close()

	log.Infof("Your identity:\n\n\t%s\n\n", local.Host.ID())

	err = local.StartMdnsService(ctx)
	if err != nil {
		return err
	}
	defer local.StopMdnsService()

	local.RegisterRequestHandler(local)

	log.Infoln("Ready to receive files... (cancel with ctrl+c)")

	return <-shutdown
}

func printInformation(data *p2p.PushRequest) {

	var cStr string
	if c, err := cid.Cast(data.Cid); err != nil {
		cStr = err.Error()
	} else {
		cStr = c.String()
	}

	log.Infoln("Sending request information:")
	log.Infoln("\tPeer:\t", data.Header.NodeId)
	log.Infoln("\tName:\t", data.Filename)
	log.Infoln("\tSize:\t", data.Size)
	log.Infoln("\tCID:\t", cStr)
	log.Infoln("\tSign:\t", hex.EncodeToString(data.Header.Signature))
	log.Infoln("\tPubKey:\t", hex.EncodeToString(data.Header.GetNodePubKey()))
}

func help() {
	log.Infoln("y: accept and thus accept the file")
	log.Infoln("n: reject the request to accept the file")
	log.Infoln("i: show information about the sender and file to be received")
	log.Infoln("q: quit p2p")
	log.Infoln("?: this help message")
}
