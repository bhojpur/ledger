package core

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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
}

func TestTransaction(t *testing.T) {
	user, err := NewUser("Tester")
	assert.NoError(t, err)

	txn, err := NewTransaction(user)
	assert.NoError(t, err)

	cash, err := NewAccount("1", "cash")
	assert.NoError(t, err)
	income, err := NewAccount("2", "income")
	assert.NoError(t, err)
	aud, err := NewCurrency("AUD", 2)
	assert.NoError(t, err)

	amountDR := big.NewInt(10)

	spl1, err := NewSplit(time.Now(), []byte("Cash Income"), []*Account{cash}, aud, amountDR)
	assert.NoError(t, err)

	err = txn.AppendSplit(spl1)
	assert.NoError(t, err)

	amountCR := big.NewInt(-10)
	spl2, err := NewSplit(time.Now(), []byte("Cash Income"), []*Account{income}, aud, amountCR)
	assert.NoError(t, err)

	err = txn.AppendSplit(spl2)
	assert.NoError(t, err)

	total, txnBalances := txn.Balance()
	assert.True(t, txnBalances)
	assert.Equal(t, total.Cmp(big.NewInt(0)), 0)

	assert.Equal(t, txn.Splits[0].Amount, amountDR)
	assert.Equal(t, txn.Splits[0].Accounts[0].Name, "cash")
	assert.Equal(t, txn.Splits[1].Amount, amountCR)
	assert.Equal(t, txn.Splits[1].Accounts[0].Name, "income")

	//func ReverseTransaction(originalTxn *Transaction, usr *User) (*Transaction, error) {
	reversedTxn, err := ReverseTransaction(txn, user)
	assert.NoError(t, err)

	total, txnBalances = reversedTxn.Balance()
	assert.True(t, txnBalances)
	assert.Equal(t, total.Cmp(big.NewInt(0)), 0)

	assert.Equal(t, reversedTxn.Splits[0].Amount, amountCR)
	assert.Equal(t, reversedTxn.Splits[0].Accounts[0].Name, "cash")
	assert.Equal(t, reversedTxn.Splits[1].Amount, amountDR)
	assert.Equal(t, reversedTxn.Splits[1].Accounts[0].Name, "income")

}
