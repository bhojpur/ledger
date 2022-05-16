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
	//"flag"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	cmd "github.com/bhojpur/ledger/cmd/server"

	"github.com/urfave/cli/v2"
)

const (
	transactionDateFormat = "2006/01/02"
	displayPrecision      = 2
)

var CommandFile = &cli.Command{
	Name:      "file",
	Usage:     "ledgerctl file ./test/transaction-codes-1.test",
	ArgsUsage: "[]",
	Description: `
	Loads a file in the Bhojpur Ledger format
`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "currency",
			Aliases: []string{"c"},
			Value:   "INR",
			Usage:   "Specify the currency that the Bhojpur Ledger file will be in, default to INR",
		},
	},
	Action: func(ctx *cli.Context) error {
		err, cfg := cmd.MakeConfig(ctx)
		if err != nil {
			return fmt.Errorf("Could not make config (%v)", err)
		}

		var ledgerFileName string

		if ctx.NArg() > 0 {
			columnWidth := 80

			ledgerFileName = "test/transaction-codes-2.test"
			if len(ctx.Args().Get(0)) > 0 {
				ledgerFileName = ctx.Args().Get(0)
			}

			ledgerFileReader, err := NewLedgerReader(ledgerFileName)
			if err != nil {
				return fmt.Errorf("Could not read file %s (%v)", ledgerFileName, err)
			}

			generalLedger, parseError := ParseLedger(ledgerFileReader, ctx.String("currency"))
			if parseError != nil {
				return fmt.Errorf("Could not parse file (%v)", parseError)
			}

			PrintLedger(generalLedger, columnWidth)
			SendLedger(cfg, generalLedger)
		} else {
			return errors.New("This command requires an argument")
		}
		return nil
	},
}

// PrintTransaction prints a transaction formatted to fit in specified column width.
func PrintTransaction(trans *Transaction, columns int) {
	fmt.Printf("%+v\n", trans)
	fmt.Printf("%s %s\n", trans.Date.Format(transactionDateFormat), trans.Payee)
	for _, accChange := range trans.AccountChanges {
		outBalanceString := accChange.Balance.FloatString(displayPrecision)
		spaceCount := columns - 4 - utf8.RuneCountInString(accChange.Name) - utf8.RuneCountInString(outBalanceString)
		fmt.Printf("    %s%s%s\n", accChange.Name, strings.Repeat(" ", spaceCount), outBalanceString)
	}
	fmt.Println("")
}

// PrintLedger prints all transactions as a formatted ledger file.
func PrintLedger(generalLedger []*Transaction, columns int) {
	for _, trans := range generalLedger {
		PrintTransaction(trans, columns)
	}
}

func SendLedger(cfg *cmd.LedgerConfig, generalLedger []*Transaction) {
	for _, trans := range generalLedger {
		Send(cfg, trans)
	}
}