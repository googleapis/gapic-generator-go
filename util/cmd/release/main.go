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
		"-n="+version,
		"-b='"+notes(version)+"'",
		"-u=googleapis",
		"-r=gapic-generator-go",
		"-c="+commitish,
		version,
		stagingDir)
}

// notes collects the commit messages since the previous tag, aggregates
// dependency updates into one message, and returns formatted release notes.
func notes(version string) string {
	previous := util.ExecuteWithOutput(
		"git",
		"describe",
		"--abbrev=0",
		"--tags",
		version+"^",
	)
	previous = strings.TrimSpace(previous)

	output := util.ExecuteWithOutput(
		"git",
		"log",
		previous+"..HEAD",
		"--oneline",
		"--pretty=format:%s",
	)

	var hasDeps bool
	var gapic, bazel, gencli, chore, samples, other []string
	for _, msg := range strings.Split(output, "\n") {
		split := strings.Split(msg, ":")
		comp := split[0]

		switch comp {
		case "gapic":
			gapic = append(gapic, "*"+split[1])
		case "bazel":
			bazel = append(bazel, "*"+split[1])
		case "gencli":
			gencli = append(gencli, "*"+split[1])
		case "chore(deps)":
			hasDeps = true
		case "chore":
			chore = append(chore, "*"+split[1])
		case "samples":
			samples = append(samples, "*"+split[1])
		default:
			other = append(other, "* "+msg)
		}
	}

	if hasDeps {
		chore = append(chore, "* update dependencies (see history)")
	}

	var notes strings.Builder
	if len(gapic) > 0 {
		notes.WriteString("# gapic\n\n")
		notes.WriteString(strings.Join(gapic, "\n"))
		notes.WriteString("\n\n")
	}

	if len(bazel) > 0 {
		notes.WriteString("# bazel\n\n")
		notes.WriteString(strings.Join(bazel, "\n"))
		notes.WriteString("\n\n")
	}

	if len(gencli) > 0 {
		notes.WriteString("# gencli\n\n")
		notes.WriteString(strings.Join(gencli, "\n"))
		notes.WriteString("\n\n")
	}

	if len(samples) > 0 {
		notes.WriteString("# samples\n\n")
		notes.WriteString(strings.Join(samples, "\n"))
		notes.WriteString("\n\n")
	}

	if len(chore) > 0 {
		notes.WriteString("# chores\n\n")
		notes.WriteString(strings.Join(chore, "\n"))
		notes.WriteString("\n\n")
	}

	if len(other) > 0 {
		notes.WriteString("# other\n\n")
		notes.WriteString(strings.Join(other, "\n"))
	}

	return notes.String()
}
