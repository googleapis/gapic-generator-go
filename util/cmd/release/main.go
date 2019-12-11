// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main contains a script that is run in CI when a new version tag is pushed to master.
// This script archives compiled executables of the gapic-generator-go plugin tool and creates
// a GitHub release with for the given tag and commitish using the given GitHub token.
//
// This script must be run from the root directory of the gapic-generator-go repository.
//
// Usage: go run ./util/cmd/release -version=v1.2.3 -token=$GITHUB_TOKEN -commitish=abc123
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/googleapis/gapic-generator-go/util"
)

var version, token, commitish string

func init() {
	flag.StringVar(&commitish, "commitish", "", "the target commitish to use [required]")
	flag.StringVar(&token, "token", "", "the GitHub token to use [optional]")
	flag.StringVar(&version, "version", "", "the version tag [required]")
}

func main() {
	flag.Parse()
	if commitish == "" {
		log.Fatalln("Missing required flag: -commitish")
	}
	if token == "" {
		log.Fatalln("Missing required flag: -token")
	}
	if version == "" {
		log.Fatalln("Missing required flag: -version")
	}
	if !strings.HasPrefix(version, "v") {
		log.Fatalln("Invalid version, should be in form vX.X.X")
	}

	outDir, err := ioutil.TempDir(os.TempDir(), "gapic-generator-go-")
	if err != nil {
		log.Fatalf("Error: unable to create temporary directory %+v\n", err)
	}
	defer os.RemoveAll(outDir)

	// Get cross compiler & GitHub release helper
	// Mousetrap is a windows dependency that is not implicitly got since
	// we only get the linux dependencies.
	util.Execute(
		"go",
		"get",
		"github.com/mitchellh/gox",
		"github.com/inconshreveable/mousetrap",
		"github.com/tcnksm/ghr")

	// Compile plugin binaries.
	stagingDir := filepath.Join(outDir, "binaries")
	osArchs := []string{
		"windows/amd64",
		"linux/amd64",
		"darwin/amd64",
		"linux/arm",
	}
	for _, osArch := range osArchs {
		util.Execute(
			"gox",
			fmt.Sprintf("-osarch=%s", osArch),
			"-output",
			filepath.Join(stagingDir, fmt.Sprintf("protoc-gen-go_gapic-%s-{{.OS}}-{{.Arch}}", version), "protoc-gen-go_gapic"),
			"github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic")
	}

	dirs, _ := filepath.Glob(filepath.Join(stagingDir, "*"))
	for _, dir := range dirs {
		// The windows binaries are suffixed with '.exe'. This allows us to create
		// tarballs of the executables whether or not they contain a suffix.
		files, _ := filepath.Glob(filepath.Join(dir, "protoc-gen-go_gapic*"))
		util.Execute(
			"tar",
			"-zcf",
			dir+".tar.gz",
			"-C",
			filepath.Dir(files[0]),
			filepath.Base(files[0]))
		// Remove the individual binary directory.
		util.Execute("rm", "-r", dir)
	}

	// Execute GitHub release of artifacts.
	util.Execute(
		"ghr",
		"-t="+token,
		"-u=googleapis",
		"-r=gapic-generator-go",
		"-c="+commitish,
		version,
		stagingDir)
}
