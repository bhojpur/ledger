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
	"flag"
	"fmt"
	"log"

	cmd "github.com/bhojpur/ledger/cmd/server"
	"github.com/bhojpur/ledger/pkg/api/v1/transaction"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

type Account struct {
	AccountName string
	Tags        []string
}

var accounts = []Account{
	{AccountName: "Cash",
		Tags: []string{"main", "Current Asset", "Asset", "Balance Sheet"}},
	{AccountName: "Accounts Receivable",
		Tags: []string{"main", "Current Asset", "Asset", "Balance Sheet"}},
	{AccountName: "Accounts Payable",
		Tags: []string{"main", "Current Liability", "Liability", "Balance Sheet"}},
	{AccountName: "Retained Earnings",
		Tags: []string{"main", "Equity", "Balance Sheet"}},
	{AccountName: "Sales",
		Tags: []string{"main", "Revenue", "Profit and Loss"}},
	{AccountName: "General Expenses",
		Tags: []string{"main", "Expense", "Profit and Loss"}},
	{AccountName: "Rent",
		Tags: []string{"main", "Expense", "Profit and Loss"}},
	{AccountName: "Interest",
		Tags: []string{"main", "Expense", "Profit and Loss"}},
	{AccountName: "Computer Expenses",
		Tags: []string{"main", "Expense", "Profit and Loss"}},
	{AccountName: "Salary and Wages",
		Tags: []string{"main", "Expense", "Profit and Loss"}},
	{AccountName: "Minor Equipment",
		Tags: []string{"main", "Expense", "Profit and Loss"}},
	{AccountName: "Repairs and Maintenance",
		Tags: []string{"main", "Expense", "Profit and Loss"}},
	{AccountName: "Bank Account",
		Tags: []string{"External"}},
}

func main() {

	// Create a config from defaults which would usually be created by Bhojpur CLI library
	set := flag.NewFlagSet("accounts", 0)
	set.String("config", "", "doc")
	ctx := cli.NewContext(nil, set, nil)
	err, config := cmd.MakeConfig(ctx)
	if err != nil {
		log.Fatalf("New Config Failed: %v", err)
	}

	conns := make([]*grpc.ClientConn, 1)
	for i := 0; i < len(conns); i++ {
		log.Printf("Starting Bhojpur Ledger %d", i)
		conn, err := grpc.Dial(fmt.Sprintf("%s:%s", config.Host, config.RPCPort), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to dial: %v", err)
		}
		conns[i] = conn
		defer func() {
			if err := conn.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	client := transaction.NewTransactorClient(conns[0])

	for i := 0; i < len(accounts); i++ {
		req := &transaction.AccountTagRequest{
			Account: accounts[i].AccountName,
			Tag:     accounts[i].Tags,
		}
		_, err = client.AddAccount(context.Background(), req)
		if err != nil {
			log.Fatal(err)
		}
	}

}
