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
	"encoding/json"
	"fmt"

	cmd "github.com/bhojpur/ledger/cmd/server"
	"github.com/urfave/cli/v2"
)

var CommandJSONJournal = &cli.Command{
	Name:      "jsonjournal",
	Usage:     "ledgerctl jsonjournal <journalInJSONFormat>",
	ArgsUsage: "[]",
	Description: `
	Creates a journal using the JSON passed through as the first Argument

	Example

	ledgerctl jsonjournal '{"Payee":"ijfjie","Date":"2019-06-30T00:00:00Z","AccountChanges":[{"Name":"Cash","Description":"jisfeij","Currency":"INR","Balance":"100"},{"Name":"Income","Description":"another","Currency":"INR","Balance":"-100"}]}'

`,
	Flags: []cli.Flag{},
	Action: func(ctx *cli.Context) error {
		err, cfg := cmd.MakeConfig(ctx)
		if err != nil {
			return fmt.Errorf("Could not make config (%v)", err)
		}

		// we initialize our request struct
		var req Transaction

		// we unmarshal our byteArray which contains our
		// jsonFile's content into 'users' which we defined above
		json.Unmarshal([]byte(ctx.Args().Get(0)), &req)

		log.Debugf("Transaction: %v\n", req)

		err = Send(cfg, &req)
		if err != nil {
			return fmt.Errorf("Could not send transaction (%v)", err)
		}

		return nil
	},
}
