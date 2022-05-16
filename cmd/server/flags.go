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

// It defines the command line flags for the shared utlities.

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v2"
)

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

// DefaultDataDir is the default data directory to use for the databases and other
// persistence requirements.
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "ledger")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "ledger")
		} else {
			return filepath.Join(home, ".bhojpur")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

var (
	// VerbosityFlag defines the logrus configuration.
	VerbosityFlag = &cli.StringFlag{
		Name:  "verbosity",
		Usage: "Logging verbosity (debug, info=default, warn, error, fatal, panic)",
	}
	// DataDirFlag defines a path on disk.
	DataDirFlag = &cli.StringFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
	}
	// ClearDB tells the node to remove any previously stored data at the data directory.
	ClearDB = &cli.BoolFlag{
		Name:  "clear-db",
		Usage: "Clears any previously stored data at the data directory",
	}
	ConfigFileFlag = &cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}
	// RPCHost defines the host on which the RPC server should listen.
	RPCHost = &cli.StringFlag{
		Name:  "rpc-host",
		Usage: "Host on which the RPC server should listen",
	}
	// RPCPort defines a beacon node RPC port to open.
	RPCPort = &cli.StringFlag{
		Name:  "rpc-port",
		Usage: "RPC port exposed by Bhojpur Ledger",
	}
	// CertFlag defines a flag for the node's TLS CA certificate.
	CACertFlag = &cli.StringFlag{
		Name:  "ca-cert",
		Usage: "Certificate Authority certificate for secure gRPC. Pass this and the tls-key flag in order to use gRPC securely.",
	}
	// CertFlag defines a flag for the node's TLS certificate.
	CertFlag = &cli.StringFlag{
		Name:  "tls-cert",
		Usage: "Certificate for secure gRPC. Pass this and the tls-key flag in order to use gRPC securely.",
	}
	// KeyFlag defines a flag for the node's TLS key.
	KeyFlag = &cli.StringFlag{
		Name:  "tls-key",
		Usage: "Key for secure gRPC. Pass this and the tls-cert flag in order to use gRPC securely.",
	}
	// LogFileName specifies the log output file name.
	LogFileName = &cli.StringFlag{
		Name:  "log-file",
		Usage: "Specify log file name, relative or absolute",
	}
	// DatabaseType specifies the backend for Bhojpur Ledger
	DatabaseTypeFlag = &cli.StringFlag{
		Name:  "database",
		Usage: "Specify database type, sqlite3 or mysql",
	}
	// DatabaseLocation specifies file location for Sqlite or connection string for MySQL
	DatabaseLocationFlag = &cli.StringFlag{
		Name:  "database-location",
		Usage: "location of database file or connection string",
	}
	PidFileFlag = &cli.StringFlag{
		Name:  "pidfile",
		Usage: "location of PID File (blank will mean none is created)",
	}
)

func setConfig(ctx *cli.Context, cfg *LedgerConfig) {

	if ctx.IsSet(VerbosityFlag.Name) {
		cfg.LogVerbosity = ctx.String(VerbosityFlag.Name)
	}
	if ctx.IsSet(ConfigFileFlag.Name) {
		cfg.ConfigFile = ctx.String(ConfigFileFlag.Name)
	}
	if ctx.IsSet(DataDirFlag.Name) {
		cfg.DataDirectory = ctx.String(DataDirFlag.Name)
	}
	if ctx.IsSet(RPCHost.Name) {
		cfg.Host = ctx.String(RPCHost.Name)
	}
	if ctx.IsSet(RPCPort.Name) {
		cfg.RPCPort = ctx.String(RPCPort.Name)
	}
	if ctx.IsSet(CACertFlag.Name) {
		cfg.CACert = ctx.String(CACertFlag.Name)
	}
	if ctx.IsSet(CertFlag.Name) {
		cfg.Cert = ctx.String(CertFlag.Name)
	}
	if ctx.IsSet(KeyFlag.Name) {
		cfg.Key = ctx.String(KeyFlag.Name)
	}
	if ctx.IsSet(DatabaseTypeFlag.Name) {
		cfg.DatabaseType = ctx.String(DatabaseTypeFlag.Name)
	}
	if ctx.IsSet(DatabaseLocationFlag.Name) {
		cfg.DatabaseLocation = ctx.String(DatabaseLocationFlag.Name)
	}
	if ctx.IsSet(PidFileFlag.Name) {
		cfg.PidFile = ctx.String(PidFileFlag.Name)
	}
}
