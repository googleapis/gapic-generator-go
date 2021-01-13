// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This tool reads from the `.git` directory and determines the upcoming
// version tag and changelog.
//
// Usage (see .github/workflows/release.yaml):
//   go run ./.github/release-tool
package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var version string

func init() {
	flag.StringVar(&version, "version", "", "the version tag [required]")
}

func main() {
	flag.Parse()
	if version == "" {
		log.Fatalln("Missing required flag: -version")
	}
	if !strings.HasPrefix(version, "v") {
		log.Fatalln("Invalid version, should be in form vX.X.X")
	}

	// Get the previously released version.
	lastVersion := mustExec("git", "describe", "--abbrev=0", "--tags", version+"^")

	// Get the changelog between the most recent release version and now.
	cl := newChangelog(
		mustExec("git", "log", fmt.Sprintf("%s..HEAD", lastVersion), "--oneline", "--pretty=format:%s"),
	)

	// Dump output.
	expr := "::set-output name=%s::%s\n"
	fmt.Printf("Previous version: %s\n", lastVersion)
	fmt.Printf("New version: %s\n", version)
	fmt.Printf(expr, "version", version)
	fmt.Printf(expr, "release_notes", cl.notes())
}

func mustExec(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		log.Fatalf("exec failed: %s\n%s", out, err)
	}
	return strings.TrimSpace(string(out))
}
