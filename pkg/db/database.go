package db

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
	"database/sql"
	"time"

	"github.com/bhojpur/ledger/pkg/core"
)

// Database wraps all database operations.
type Database interface {
	InitDB() error
	Close() error
	AddTransaction(txn *core.Transaction) (string, error)
	FindTransaction(txnID string) (*core.Transaction, error)
	DeleteTransaction(txnID string) error
	FindTag(tag string) (int, error)
	AddTag(tag string) error
	SafeAddTag(tag string) error
	SafeAddTagToAccount(account, tag string) error
	AddTagToAccount(accountID string, tag int) error
	DeleteTagFromAccount(account, tag string) error
	SafeAddTagToTransaction(txnID, tag string) error
	AddTagToTransaction(txnID string, tag int) error
	DeleteTagFromTransaction(txnID, tag string) error
	FindCurrency(cur string) (*core.Currency, error)
	AddCurrency(cur *core.Currency) error
	SafeAddCurrency(cur *core.Currency) error
	DeleteCurrency(currency string) error
	FindAccount(code string) (*core.Account, error)
	AddAccount(*core.Account) error
	SafeAddAccount(*core.Account) (bool, error)
	DeleteAccount(accountName string) error
	FindUser(pubKey string) (*core.User, error)
	AddUser(usr *core.User) error
	ReconcileTransactions(reconciliationID string, splitIDs []string) (string, error)
	SafeAddUser(usr *core.User) error
	GetTB(date time.Time) (*[]core.TBAccount, error)
	GetListing(startdate, enddate time.Time) (*[]core.Transaction, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}
