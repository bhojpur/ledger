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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	cmd "github.com/bhojpur/ledger/cmd/server"
	"github.com/bhojpur/ledger/pkg/api/v1/transaction"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/urfave/cli/v2"
)

var CommandSingleTestTransaction = &cli.Command{
	Name:      "single",
	Usage:     "submits a single transaction",
	ArgsUsage: "[]",
	Description: `
`,
	Flags: []cli.Flag{},
	Action: func(ctx *cli.Context) error {
		err, cfg := cmd.MakeConfig(ctx)
		if err != nil {
			return fmt.Errorf("Could not make config (%v)", err)
		}

		date, _ := time.Parse("2006-01-02", "2011-03-15")
		desc := "Whole Food Market"

		transactionLines := make([]Account, 2)

		line1Account := "Expenses:Groceries"
		line1Desc := "Groceries"
		line1Amount := big.NewRat(7500, 1)

		transactionLines[0] = Account{
			Name:        line1Account,
			Description: line1Desc,
			Balance:     line1Amount,
			Currency:    "INR",
		}

		line2Account := "Assets:Checking"
		line2Desc := "Groceries"
		line2Amount := big.NewRat(-7500, 1)

		transactionLines[1] = Account{
			Name:        line2Account,
			Description: line2Desc,
			Balance:     line2Amount,
			Currency:    "INR",
		}

		req := &Transaction{
			Date:           date,
			Payee:          desc,
			AccountChanges: transactionLines,
		}

		err = Send(cfg, req)
		if err != nil {
			return fmt.Errorf("Could not send transaction (%v)", err)
		}

		return nil
	},
}

func Send(cfg *cmd.LedgerConfig, t *Transaction) error {

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
	//conn, err := grpc.Dial(address, grpc.WithInsecure())
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return fmt.Errorf("Could not connect to gRPC (%v)", err)
	}
	defer conn.Close()
	client := transaction.NewTransactorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	transactionLines := make([]*transaction.LineItem, 2)

	for i, accChange := range t.AccountChanges {
		amountInt64 := accChange.Balance.Num().Int64() * int64(100) / accChange.Balance.Denom().Int64()
		transactionLines[i] = &transaction.LineItem{
			Accountname: accChange.Name,
			Description: accChange.Description,
			Amount:      amountInt64,
			Currency:    accChange.Currency,
		}
	}

	req := &transaction.TransactionRequest{
		Date:        t.Date.Format("2006-01-02"),
		Description: t.Payee,
		Lines:       transactionLines,
	}
	r, err := client.AddTransaction(ctx, req)
	if err != nil {
		return fmt.Errorf("Could not call Add Transaction Method (%v)", err)
	}
	log.Infof("Add Transaction Response: %s", r.GetMessage())
	return nil
}
func loadTLSCredentials(cfg *cmd.LedgerConfig) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile(cfg.CACert)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}
