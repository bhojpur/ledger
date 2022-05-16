package version

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
	"fmt"
	"strings"
)

const (
	VersionMajor = 0       // Major version component of the current release
	VersionMinor = 0       // Minor version component of the current release
	VersionPatch = 1       // Patch version component of the current release
	VersionMeta  = "alpha" // Version metadata to append to the version string
)

var (
	// Version is the semver release name of this build
	Version string = func() string {
		return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	}()
	// Commit is the commit hash this build was created from
	Commit string
	// Branch is the branch hash this build was created from
	Branch string
	// Date is the time when this build was created
	Date string

	// GitCommit will be overwritten automatically by the build system
	BuildTime string
	// BuildCommit will be overwritten automatically by the build system
	BuildCommit = "HEAD"
)

// Print writes the version info to stdout
func Print() {
	fmt.Printf("Version:    %s\n", Version)
	fmt.Printf("Commit:     %s\n", Commit)
	fmt.Printf("Build Date: %s\n", Date)
}

// FullVersion formats the version to be printed
func FullVersion() string {
	return fmt.Sprintf("%s (%s, build %s)", Version, BuildTime, BuildCommit)
}

// RC checks if the Bhojpur Ledger version is a release candidate or not
func RC() bool {
	return strings.Contains(Version, "rc")
}

// VersionWithMeta holds the textual version string including the metadata.
var VersionWithMeta = func() string {
	v := Version
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}()

// ArchiveVersion holds the textual version string used for Swarm archives.
// e.g. "0.3.0-dea1ce05" for stable releases, or
//      "0.3.1-unstable-21c059b6" for unstable releases
func ArchiveVersion(gitCommit string) string {
	vsn := Version
	if VersionMeta != "stable" {
		vsn += "-" + VersionMeta
	}
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	return vsn
}

func VersionWithCommit() string {
	vsn := Version
	vsn += "-" + Branch
	if len(Commit) >= 8 {
		vsn += "-" + Commit[:8]
	}
	vsn += "-" + Date
	return vsn
}
