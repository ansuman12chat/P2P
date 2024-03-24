package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/ansuman12chat/p2p/internal/log"
	"github.com/ansuman12chat/p2p/pkg/receive"
	"github.com/ansuman12chat/p2p/pkg/send"
)

var (
	RawVersion = "1.0.0"
)

func main() {
	verTag := fmt.Sprintf("v%s", RawVersion)
	app := &cli.App{
		Name: "p2p",
		Authors: []*cli.Author{
			{
				Name:  "Ansuman Singh",
				Email: "ansuman12chat@gmail.com",
			},
		},
		Usage:                "A peer-to-peer data transfer tool.",
		Version:              verTag,
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			send.Command,
			receive.Command,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Infof("error: %v\n", err)
		os.Exit(1)
	}
}
