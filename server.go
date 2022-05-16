//go:build !server
// +build !server

package main

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	cmd "github.com/bhojpur/ledger/cmd/server"
	engine "github.com/bhojpur/ledger/pkg/engine"
	"github.com/bhojpur/ledger/pkg/node"
	"github.com/bhojpur/ledger/pkg/rpc"
	"github.com/bhojpur/ledger/pkg/version"
)

func startNode(ctx *cli.Context) error {
	err, cfg := cmd.MakeConfig(ctx)
	if err != nil {
		return err
	}

	fullnode, err := node.New(ctx, cfg)
	if err != nil {
		return err
	}
	ledger, err := engine.New(ctx, cfg)
	if err != nil {
		return err
	}
	fullnode.Register(ledger)
	rpc := rpc.NewRPCService(context.Background(), &rpc.Config{
		Host:       cfg.Host,
		Port:       cfg.RPCPort,
		CACertFlag: cfg.CACert,
		CertFlag:   cfg.Cert,
		KeyFlag:    cfg.Key,
	}, ledger)
	fullnode.Register(rpc)
	fullnode.Start()

	return nil
}

func main() {
	customFormatter := new(prefixed.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
	log := logrus.WithField("prefix", "main")
	app := cli.NewApp()
	app.Name = "ledgersvr"
	app.Usage = "A financial accounting system, double-entry records server engine"
	app.EnableBashCompletion = true
	app.Action = startNode
	app.Before = func(ctx *cli.Context) error {
		// If persistent log files are written - we disable the log messages coloring because
		// the colors are ANSI codes and seen as Gibberish in the log files.
		logFileName := ctx.String(cmd.LogFileName.Name)
		if logFileName != "" {
			customFormatter.DisableColors = true
			if err := cmd.ConfigurePersistentLogging(logFileName); err != nil {
				log.WithError(err).Error("Failed to configuring logging to disk.")
			}
		}
		return nil
	}
	app.Version = version.VersionWithCommit()
	app.Commands = []*cli.Command{
		// See cmd/config.go
		cmd.DumpConfigCommand,
		cmd.GenConfigCommand,
	}

	app.Flags = []cli.Flag{
		// See cmd/flags.go
		cmd.VerbosityFlag,
		cmd.DataDirFlag,
		cmd.ClearDB,
		cmd.ConfigFileFlag,
		cmd.RPCHost,
		cmd.RPCPort,
		cmd.CACertFlag,
		cmd.CertFlag,
		cmd.KeyFlag,
		cmd.LogFileName,
		cmd.DatabaseTypeFlag,
		cmd.DatabaseLocationFlag,
		cmd.PidFileFlag,
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
