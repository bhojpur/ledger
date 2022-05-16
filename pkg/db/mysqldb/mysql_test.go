package mysqldb

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
	//"database/sql"
	//"log"
	"testing"

	//"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestValidateConnectionString(t *testing.T) {

	//Test Regular string with no params
	validatedString, err := ValidateConnectionString("ledger:password@tcp(192.168.1.98:3306)/ledger")
	assert.Nil(t, err)
	if assert.NotNil(t, validatedString, "Connection String with no params") {
		assert.Equal(t, "ledger:password@tcp(192.168.1.98:3306)/ledger?parseTime=true&charset=utf8", validatedString)
	}

	//Test Same String with a param to ensure no duplication
	validatedString, err = ValidateConnectionString("ledger:password@tcp(192.168.1.98:3306)/ledger?parseTime=true")
	assert.Nil(t, err)
	if assert.NotNil(t, validatedString, "Connection String with param") {
		assert.Equal(t, "ledger:password@tcp(192.168.1.98:3306)/ledger?parseTime=true&charset=utf8", validatedString)
	}

	//Test empty string to ensure error
	nilString, err := ValidateConnectionString("")
	if assert.NotNil(t, err, "Nil connection string") {
		assert.Equal(t, "Connection string not provided", err.Error())
	}
	assert.Equal(t, "", nilString)
}
