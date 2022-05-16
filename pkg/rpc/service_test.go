package rpc

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
	"errors"
	"flag"
	"io/ioutil"
	"testing"

	cmd "github.com/bhojpur/ledger/cmd/server"
	engine "github.com/bhojpur/ledger/pkg/engine"
	"github.com/bhojpur/ledger/pkg/internal"

	"github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(ioutil.Discard)
}

func TestLifecycle_OK(t *testing.T) {
	hook := logTest.NewGlobal()

	set := flag.NewFlagSet("test", 0)
	set.String("config", "", "doc")
	ctx := cli.NewContext(nil, set, nil)
	err, cfg := cmd.MakeConfig(ctx)
	assert.NoError(t, err)

	cfg.DatabaseType = "memorydb"
	cfg.Host = "127.0.0.1"
	cfg.RPCPort = "7348"
	cfg.CACert = "bob.crt"
	cfg.Cert = "alice.crt"
	cfg.Key = "alice.key"

	ledger, err := engine.New(ctx, cfg)
	assert.NoError(t, err)

	rpcService := NewRPCService(context.Background(), &Config{
		Host:       cfg.Host,
		Port:       cfg.RPCPort,
		CACertFlag: cfg.CACert,
		CertFlag:   cfg.Cert,
		KeyFlag:    cfg.Key,
	}, ledger)

	rpcService.Start()

	internal.LogsContain(t.Fatalf, hook, "gRPC Listening on port", true)
	assert.NoError(t, rpcService.Stop())
}

func TestStatus_CredentialError(t *testing.T) {
	credentialErr := errors.New("credentialError")
	s := &Service{credentialError: credentialErr}

	assert.Contains(t, s.credentialError.Error(), s.Status().Error())
}

func TestRPC_InsecureEndpoint(t *testing.T) {
	hook := logTest.NewGlobal()

	set := flag.NewFlagSet("test", 0)
	set.String("config", "", "doc")
	ctx := cli.NewContext(nil, set, nil)
	err, cfg := cmd.MakeConfig(ctx)
	assert.NoError(t, err)

	cfg.DatabaseType = "memorydb"
	cfg.Host = "127.0.0.1"
	cfg.RPCPort = "7777"

	ledger, err := engine.New(ctx, cfg)
	assert.NoError(t, err)

	rpcService := NewRPCService(context.Background(), &Config{
		Host: cfg.Host,
		Port: cfg.RPCPort,
	}, ledger)

	rpcService.Start()

	internal.LogsContain(t.Fatalf, hook, "gRPC Listening on port", true)
	internal.LogsContain(t.Fatalf, hook, "You are using an insecure gRPC server", true)
	assert.NoError(t, rpcService.Stop())
}
