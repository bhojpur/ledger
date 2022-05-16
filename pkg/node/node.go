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
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	cmd "github.com/bhojpur/ledger/cmd/server"
	"github.com/bhojpur/ledger/pkg/core"
	"github.com/bhojpur/ledger/pkg/db"
	"github.com/bhojpur/ledger/pkg/version"
)

var log = logrus.WithField("prefix", "node")

type Node struct {
	ctx      *cli.Context
	lock     sync.RWMutex
	services *core.ServiceRegistry
	stop     chan struct{} // Channel to wait for termination notifications.
	DB       *db.Database
	PidFile  string
}

func New(ctx *cli.Context, cfg *cmd.LedgerConfig) (*Node, error) {

	registry := core.NewServiceRegistry()

	ledger := &Node{
		ctx:      ctx,
		services: registry,
		stop:     make(chan struct{}),
		PidFile:  cfg.PidFile,
	}

	return ledger, nil

}

func (n *Node) Register(constructor core.Service) error {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.services.RegisterService(constructor)

	return nil
}

func (n *Node) Start() {
	n.lock.Lock()
	log.WithFields(logrus.Fields{
		"version": version.VersionWithCommit(),
	}).Info("Starting Bhojpur Ledger server")

	n.writePIDFile()
	n.services.StartAll()
	stop := n.stop
	n.lock.Unlock()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		<-sigc
		log.Info("Got interrupt, shutting down...")
		go n.Close()
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				log.Info("Already shutting down, interrupt more to panic", "times", i-1)
			}
		}
		panic("Panic closing the Bhojpur Ledger node")
	}()

	<-stop
}

// Close handles graceful shutdown of the system.
func (n *Node) Close() {
	n.lock.Lock()
	defer n.lock.Unlock()

	log.Info("Stopping Bhojpur Ledger node")
	if len(n.PidFile) > 0 {
		_ = os.Remove(n.PidFile)
	}
	n.services.StopAll()
	close(n.stop)
}

// writePIDFile retrieves the current process ID and writes it to file.
func (n *Node) writePIDFile() {

	if n.PidFile == "" {
		return
	}

	// Ensure the required directory structure exists.
	err := os.MkdirAll(filepath.Dir(n.PidFile), 0700)
	if err != nil {
		log.Error("Failed to verify pid directory", "error", err)
		os.Exit(1)
	}

	// Retrieve the PID and write it to file.
	pid := strconv.Itoa(os.Getpid())
	if err := ioutil.WriteFile(n.PidFile, []byte(pid), 0644); err != nil {
		log.Error("Failed to write pidfile", "error", err)
		os.Exit(1)
	}

	log.Info("Writing PID file", "path", n.PidFile, "pid", pid)
}
