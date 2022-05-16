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
	"fmt"
	"os"

	clicmd "github.com/bhojpur/ledger/cmd/reporter/commands"
	cmd "github.com/bhojpur/ledger/cmd/server"
	"github.com/urfave/cli/v2"
)

var app *cli.App

func init() {
	clicmd.InitLogger()
	app = cli.NewApp()
	app.Name = "ledgerepo"
	app.Usage = "Extracts General Ledger and Trial Balance reports from a database"
	app.Commands = []*cli.Command{
		// transactionlisting.go
		clicmd.CommandTransactionListing,
		// trialbalance.go
		clicmd.CommandTrialBalance,
		// pdfgenerator.go
		clicmd.CommandPDFGenerate,
	}
	app.Flags = []cli.Flag{
		cmd.VerbosityFlag,
		cmd.ConfigFileFlag,
		cmd.RPCHost,
		cmd.RPCPort,
		cmd.CertFlag,
		cmd.KeyFlag,
	}
	app.Action = clicmd.ReporterConsole
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
