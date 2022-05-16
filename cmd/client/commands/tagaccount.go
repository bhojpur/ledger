package cmd

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
	"fmt"
	"time"

	cmd "github.com/bhojpur/ledger/cmd/server"
	"github.com/bhojpur/ledger/pkg/api/v1/transaction"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var CommandTagAccount = &cli.Command{
	Name:      "tag",
	Usage:     "ledgerctl tag <account> <tag>",
	ArgsUsage: "[]",
	Description: `
	Adds the tag specified in the second argument to the account specified in the first argument
`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "deletes tag rather than creates",
		},
	},
	Action: func(ctx *cli.Context) error {
		err, cfg := cmd.MakeConfig(ctx)
		if err != nil {
			return fmt.Errorf("Could not make config (%v)", err)
		}

		address := fmt.Sprintf("%s:%s", cfg.Host, cfg.RPCPort)
		log.WithField("address", address).Info("gRPC dialing on port")
		opts := []grpc.DialOption{}

		if cfg.CACert != "" && cfg.Cert != "" && cfg.Key != "" {
			tlsCredentials, err := loadTLSCredentials(cfg)
			if err != nil {
				return fmt.Errorf("Could not load TLS credentials (%v)", err)
			}
			opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
		} else {
			opts = append(opts, grpc.WithInsecure())
		}

		// Set up a connection to the server.
		conn, err := grpc.Dial(address, opts...)
		if err != nil {
			return fmt.Errorf("Could not connect to gRPC (%v)", err)
		}
		defer conn.Close()
		client := transaction.NewTransactorClient(conn)

		ctxtimeout, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if ctx.Bool("delete") {
			req := &transaction.DeleteAccountTagRequest{
				Account: ctx.Args().Get(0),
				Tag:     []string{ctx.Args().Get(1)},
			}

			r, err := client.DeleteTag(ctxtimeout, req)
			if err != nil {
				return fmt.Errorf("Could not call Delete Tag Method (%v)", err)
			}

			log.Infof("Delete Tag Response: %s", r.GetMessage())
		} else {
			req := &transaction.AccountTagRequest{
				Account: ctx.Args().Get(0),
				Tag:     []string{ctx.Args().Get(1)},
			}

			r, err := client.AddTag(ctxtimeout, req)
			if err != nil {
				return fmt.Errorf("Could not call Add Tag Method (%v)", err)
			}

			log.Infof("Create Tag Response: %s", r.GetMessage())
		}

		return nil
	},
}
