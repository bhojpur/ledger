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
	"math/big"
	"time"
)

// Account holds the name and balance
type Account struct {
	Name        string
	Description string
	Currency    string
	Balance     *big.Rat
}

type sortAccounts []*Account

func (s sortAccounts) Len() int      { return len(s) }
func (s sortAccounts) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type sortAccountsByName struct{ sortAccounts }

func (s sortAccountsByName) Less(i, j int) bool {
	return s.sortAccounts[i].Name < s.sortAccounts[j].Name
}

// Transaction is the basis of a ledger. The ledger holds a list of transactions.
// A Transaction has a Payee, Date (with no time, or to put another way, with
// hours,minutes,seconds values that probably doesn't make sense), and a list of
// Account values that hold the value of the transaction for each account.
type Transaction struct {
	Payee          string
	Date           time.Time
	AccountChanges []Account
}

type sortTransactions []*Transaction

func (s sortTransactions) Len() int      { return len(s) }
func (s sortTransactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type sortTransactionsByDate struct{ sortTransactions }

func (s sortTransactionsByDate) Less(i, j int) bool {
	return s.sortTransactions[i].Date.Before(s.sortTransactions[j].Date)
}