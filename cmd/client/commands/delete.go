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
	"errors"
	"fmt"
	"time"

	cmd "github.com/bhojpur/ledger/cmd/server"
	"github.com/bhojpur/ledger/pkg/api/v1/transaction"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var CommandDeleteTransaction = &cli.Command{
	Name:      "delete",
	Usage:     "ledgerctl delete <transaction_id>",
	ArgsUsage: "[]",
	Description: `
	Deletes a transaction from the database
`,
	Flags: []cli.Flag{},
	Action: func(ctx *cli.Context) error {
		err, cfg := cmd.MakeConfig(ctx)
		if err != nil {
			return fmt.Errorf("Could not make config (%v)", err)
		}

		if ctx.NArg() > 0 {
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

			req := &transaction.DeleteRequest{
				Identifier: ctx.Args().Get(0),
			}
			r, err := client.DeleteTransaction(ctxtimeout, req)
			if err != nil {
				return fmt.Errorf("Could not call Delete Transaction Method (%v)", err)
			}
			log.Infof("Delete Transaction Response: %s", r.GetMessage())
		} else {
			return errors.New("This command requires an argument")
		}

		return nil
	},
}

var CommandVoidTransaction = &cli.Command{
	Name:      "void",
	Usage:     "ledgerctl void <transaction_id>",
	ArgsUsage: "[]",
	Description: `
	Reverses a transaction by creating an identical inverse and tags both transactions as void 
`,
	Flags: []cli.Flag{},
	Action: func(ctx *cli.Context) error {
		err, cfg := cmd.MakeConfig(ctx)
		if err != nil {
			return fmt.Errorf("Could not make config (%v)", err)
		}

		if ctx.NArg() > 0 {
			address := fmt.Sprintf("%s:%s", cfg.Host, cfg.RPCPort)
			log.WithField("address", address).Info("gRPC Dialing on port")
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

			req := &transaction.DeleteRequest{
				Identifier: ctx.Args().Get(0),
			}
			r, err := client.VoidTransaction(ctxtimeout, req)
			if err != nil {
				return fmt.Errorf("Could not call Void Transaction Method (%v)", err)
			}
			log.Infof("Void Transaction Response: %s", r.GetMessage())
		} else {
			return errors.New("This command requires an argument")
		}

		return nil
	},
}
