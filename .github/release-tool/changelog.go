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

package main

import (
	"fmt"
	"strings"
)

type commit struct {
	category string
	message  string
}

func newCommit(line string) *commit {
	var cat, msg string
	lineSplit := strings.SplitN(line, ":", 2)
	if len(lineSplit) > 1 {
		cat = strings.TrimSpace(lineSplit[0])
		msg = strings.TrimSpace(lineSplit[1])
	} else {
		cat = "other"
		msg = strings.TrimSpace(lineSplit[0])
	}
	return &commit{category: cat, message: msg}
}

func (c *commit) String() string {
	return c.message
}

type changelog struct {
	gapic  []*commit
	bazel  []*commit
	gencli []*commit
	chore  []*commit
	other  []*commit
}

func newChangelog(gitlog string) *changelog {
	var gapic, bazel, gencli, chore, other []*commit
	var hasDeps bool
	for _, line := range strings.Split(gitlog, "\n") {
		cmt := newCommit(line)
		switch cmt.category {
		case "gapic":
			gapic = append(gapic, cmt)
		case "bazel":
			bazel = append(bazel, cmt)
		case "gencli":
			gencli = append(gencli, cmt)
		case "fix(deps)":
			fallthrough
		case "chore(deps)":
			hasDeps = true
		case "chore":
			chore = append(chore, cmt)
		default:
			other = append(other, cmt)
		}
	}

	if hasDeps {
		chore = append(chore, &commit{category: "chore", message: "update dependencies (see history)"})
	}

	return &changelog{
		gapic:  gapic,
		bazel:  bazel,
		gencli: gencli,
		chore:  chore,
		other:  other,
	}
}

func (cl *changelog) notes() string {
	section := func(title string, commits []*commit) string {
		if len(commits) > 0 {
			// %0A is newline: https://github.community/t/set-output-truncates-multiline-strings/16852
			answer := fmt.Sprintf("## %s%%0A%%0A", title)
			for _, cmt := range commits {
				answer += fmt.Sprintf("* %s%%0A", cmt)
			}
			return answer + "%0A"
		}
		return ""
	}
	return strings.TrimSuffix(
		section("gapic", cl.gapic)+section("bazel", cl.bazel)+section("gencli", cl.gencli)+section("chore", cl.chore)+section("other", cl.other),
		"%0A",
	)
}
