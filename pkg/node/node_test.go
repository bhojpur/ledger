package node

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
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	cmd "github.com/bhojpur/ledger/cmd/server"
	"github.com/bhojpur/ledger/pkg/internal"
	engine "github.com/bhojpur/ledger/pkg/engine"

	"github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/urfave/cli/v2"
)

const (
	maxPollingWaitTime = 1 * time.Second
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(ioutil.Discard)
}

// Test that Bhojpur Ledger node can close.
func TestNodeClose_OK(t *testing.T) {
	hook := logTest.NewGlobal()

	app := cli.App{}
	set := flag.NewFlagSet("test", 0)
	set.String("config", "", "doc")

	ctx := cli.NewContext(&app, set, nil)

	err, cfg := cmd.MakeConfig(ctx)
	assert.NoError(t, err)

	node, err := New(ctx, cfg)
	assert.NoError(t, err)

	node.Close()

	internal.LogsContain(t.Fatalf, hook, "Stopping ledger node", true)
}

// TestClearDB tests clearing the database
func TestClearDB(t *testing.T) {
	hook := logTest.NewGlobal()

	randPath, err := rand.Int(rand.Reader, big.NewInt(1000000))
	assert.NoError(t, err, "Could not generate random number for file path")
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("datadirtest%d", randPath))
	assert.NoError(t, os.RemoveAll(tmp))

	app := cli.App{}
	set := flag.NewFlagSet("test", 0)
	set.Bool(cmd.ClearDB.Name, true, "")

	ctx := cli.NewContext(&app, set, nil)
	assert.NoError(t, err)
	err, cfg := cmd.MakeConfig(ctx)
	assert.NoError(t, err)
	cfg.DatabaseType = "memorydb"
	cfg.DataDirectory = tmp

	node, err := New(ctx, cfg)
	assert.NoError(t, err)

	ledger, err := engine.New(ctx, cfg)
	assert.NoError(t, err)

	node.Register(ledger)
	go node.Start()
	d := time.Now().Add(maxPollingWaitTime)
	contextWithDeadline, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	<-contextWithDeadline.Done()
	internal.LogsContain(t.Fatalf, hook, "Clearing SQLite3 DB", true)
	//case <-contextWithDeadline.Done():

	node.Close()
	assert.NoError(t, os.RemoveAll(tmp))
}
