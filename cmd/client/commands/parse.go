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
	"bufio"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"sort"
	"strings"

	date "github.com/joyt/godate"
	"github.com/marcmak/calc/calc"
)

const (
	whitespace = " \t"
)

// ParseLedger parses a ledger file and returns a list of Transactions.
//
// Transactions are sorted by date.
func ParseLedger(ledgerReader io.Reader, currency string) (generalLedger []*Transaction, err error) {
	parseLedger(ledgerReader, currency, func(t *Transaction, e error) (stop bool) {
		if e != nil {
			err = e
			stop = true
			return
		}

		generalLedger = append(generalLedger, t)
		return
	})

	if err != nil {
		sort.Sort(sortTransactionsByDate{generalLedger})
	}

	return
}

var accountToAmountSpace = regexp.MustCompile(" {2,}|\t+")

func parseLedger(ledgerReader io.Reader, currency string, callback func(t *Transaction, err error) (stop bool)) {
	var trans *Transaction
	scanner := bufio.NewScanner(ledgerReader)
	var line string
	var filename string
	var lineCount int

	errorMsg := func(msg string) (stop bool) {
		return callback(nil, fmt.Errorf("%s:%d: %s", filename, lineCount, msg))
	}

	for scanner.Scan() {
		line = scanner.Text()

		// update filename/line if sentinel comment is found
		if strings.HasPrefix(line, markerPrefix) {
			filename, lineCount = parseMarker(line)
			continue
		}

		// remove heading and tailing space from the line
		trimmedLine := strings.Trim(line, whitespace)
		lineCount++

		// handle comments
		if commentIdx := strings.Index(trimmedLine, ";"); commentIdx >= 0 {
			trimmedLine = trimmedLine[:commentIdx]
			if len(trimmedLine) == 0 {
				continue
			}
		}

		if len(trimmedLine) == 0 {
			if trans != nil {
				transErr := balanceTransaction(trans)
				if transErr != nil {
					errorMsg("Unable to balance transaction, " + transErr.Error())
				}
				callback(trans, nil)
				trans = nil
			}
		} else if trans == nil {
			lineSplit := strings.SplitN(line, " ", 2)
			if len(lineSplit) != 2 {
				if errorMsg("Unable to parse payee line: " + line) {
					return
				}
				continue
			}
			dateString := lineSplit[0]
			transDate, dateErr := date.Parse(dateString)
			if dateErr != nil {
				errorMsg("Unable to parse date: " + dateString)
			}
			payeeString := lineSplit[1]
			trans = &Transaction{Payee: payeeString, Date: transDate}
		} else {
			var accChange Account
			lineSplit := accountToAmountSpace.Split(trimmedLine, -1)
			var nonEmptyWords []string
			for _, word := range lineSplit {
				if len(word) > 0 {
					nonEmptyWords = append(nonEmptyWords, word)
				}
			}
			lastIndex := len(nonEmptyWords) - 1
			balErr, rationalNum := getBalance(strings.Trim(nonEmptyWords[lastIndex], whitespace))
			accChange.Currency = currency
			if !balErr {
				// Assuming no balance and whole line is account name
				accChange.Name = strings.Join(nonEmptyWords, " ")
			} else {
				accChange.Name = strings.Join(nonEmptyWords[:lastIndex], " ")
				accChange.Balance = rationalNum
			}
			trans.AccountChanges = append(trans.AccountChanges, accChange)
		}
	}
	// If the file does not end on empty line, we must attempt to balance last
	// transaction of the file.
	if trans != nil {
		transErr := balanceTransaction(trans)
		if transErr != nil {
			errorMsg("Unable to balance transaction, " + transErr.Error())
		}
		callback(trans, nil)
	}
}

func getBalance(balance string) (bool, *big.Rat) {
	rationalNum := new(big.Rat)
	if strings.Contains(balance, "(") {
		rationalNum.SetFloat64(calc.Solve(balance))
		return true, rationalNum
	}
	_, isValid := rationalNum.SetString(balance)
	return isValid, rationalNum
}

// Takes a transaction and balances it. This is mainly to fill in the empty part
// with the remaining balance.
func balanceTransaction(input *Transaction) error {
	balance := new(big.Rat)
	var emptyFound bool
	var emptyAccIndex int
	for accIndex, accChange := range input.AccountChanges {
		if accChange.Balance == nil {
			if emptyFound {
				return fmt.Errorf("more than one account change empty")
			}
			emptyAccIndex = accIndex
			emptyFound = true
		} else {
			balance = balance.Add(balance, accChange.Balance)
		}
	}
	if balance.Sign() != 0 {
		if !emptyFound {
			return fmt.Errorf("no empty account change to place extra balance")
		}
	}
	if emptyFound {
		input.AccountChanges[emptyAccIndex].Balance = balance.Neg(balance)
	}
	return nil
}
