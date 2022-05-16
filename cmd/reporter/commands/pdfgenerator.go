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
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"

	cmd "github.com/bhojpur/ledger/cmd/server"
	engine "github.com/bhojpur/ledger/pkg/engine"
)

type Tag struct {
	Name     string       `json:"Name"`
	Total    int          `json:"Total"`
	Accounts []PDFAccount `json:"Accounts"`
}

type PDFAccount struct {
	Account string `json:"Account"`
	Amount  int    `json:"Amount"`
}

var reporteroutput struct {
	Data      []Tag `json:"Tags"`
	Profit    int   `json:"Profit"`
	NetAssets int   `json:"NetAssets"`
}

var CommandPDFGenerate = &cli.Command{
	Name:  "pdf",
	Usage: "ledgerepo pdf",
	Description: `
Creates a .pdf report using node.js and handlebars templates

Downloads a premade handlebars template and creates reports using the tagged accounts.

requires Nodejs on the machine and also handlebars (npm install -g handlebars) and puppeteer 
`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "template, t",
			Value: "profitandloss",
			Usage: "The name of the HTML template to create a PDF of",
		},
	},
	Action: func(ctx *cli.Context) error {

		//Check if keyfile path given and make sure it doesn't already exist.
		err, cfg := cmd.MakeConfig(ctx)
		if err != nil {
			return fmt.Errorf("Could not make config (%v)", err)
		}
		ledger, err := engine.New(ctx, cfg)
		if err != nil {
			return fmt.Errorf("Could not make new Bhojpur Ledger (%v)", err)
		}

		queryDB := `
			SELECT
					tags.tag_name,
					Table_Aggregate.account_id,
					sums
			FROM account_tag
					join ((SELECT
							split_accounts.account_id as account_id,
							SUM(splits.amount) as sums
						FROM splits 
							JOIN split_accounts 
							ON splits.split_id = split_accounts.split_id
						GROUP  BY split_accounts.account_id
							
						)) AS Table_Aggregate
						on account_tag.account_id = Table_Aggregate.account_id
					join tags
						on tags.tag_id = account_tag.tag_id
			order BY tags.tag_name
		;`

		log.Debugf("Quering the Database")
		rows, err := ledger.LedgerDb.Query(queryDB)
		if err != nil {
			return fmt.Errorf("Could not query database (%v)", err)
		}
		defer rows.Close()
		accounts := make(map[string][]PDFAccount)
		totals := make(map[string]int)

		for rows.Next() {
			var t PDFAccount
			var name string
			if err := rows.Scan(&name, &t.Account, &t.Amount); err != nil {
				return fmt.Errorf("Could not scan rows of query (%v)", err)
			}
			log.Debugf("%v", t)
			if val, ok := accounts[name]; ok {
				accounts[name] = append(val, t)
				totals[name] = totals[name] + t.Amount
			} else {
				accounts[name] = []PDFAccount{t}
				totals[name] = t.Amount
			}
		}
		if rows.Err() != nil {
			return fmt.Errorf("rows errored with (%v)", rows.Err())
		}

		for k, v := range accounts {
			reporteroutput.Data = append(reporteroutput.Data, Tag{k, totals[k], v})
		}

		dir := "src"

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return fmt.Errorf("Making Directory %s failed (%v)", dir, err)
			}
		}

		outputJson, _ := json.Marshal(reporteroutput)
		err = ioutil.WriteFile("src/output.json", outputJson, 0644)
		if err != nil {
			return fmt.Errorf("writing output.json failed (%v)", err)
		}

		if err := DownloadFile("./src/pdfgenerator.js", "https://raw.githubusercontent.com/bhojpur/ledger/master/pdfgenerator.js"); err != nil {
			return fmt.Errorf("downloading pdfgenerator.js failed (%v)", err)
		}

		filename := "./src/financials.html"
		//httpfile := "https://raw.githubusercontent.com/bhojpur/ledger/master/financials.html"
		httpfile := "https://raw.githubusercontent.com/bhojpur/ledger/master/templates/" + ctx.String("template") + ".html"

		log.Debugf("Downloading template: %s", httpfile)
		if err := DownloadFile(filename, httpfile); err != nil {
			return fmt.Errorf("downloading template %s failed (%v)", httpfile, err)
		}

		command := "node ./pdfgenerator.js"
		parts := strings.Fields(command)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Dir = "./src"

		cmd.Run()

		//Restructure and Cleanup
		err = os.Rename("src/mypdf.pdf", ctx.String("template")+".pdf")
		if err != nil {
			return fmt.Errorf("renaming file failed (%v)", err)
		}
		err = os.RemoveAll("src")
		if err != nil {
			return fmt.Errorf("removing src folder failed (%v)", err)
		}

		return nil
	},
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("downloading %s failed (%v)", url, err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("creating file %s failed (%v)", filepath, err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
