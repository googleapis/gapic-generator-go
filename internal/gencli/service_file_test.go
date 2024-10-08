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

package gencli

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"

	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
)

func TestServiceFile(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Generating the service_file panicked: %v", r)
		}
	}()

	name := "Todo"

	g := &gcli{
		root:   "Root",
		format: true,
		imports: map[string]*pbinfo.ImportSpec{
			"test": &pbinfo.ImportSpec{Name: "proto", Path: "github.com/golang/protobuf/proto"},
		},
		subcommands: map[string][]*Command{
			name: []*Command{
				&Command{MethodCmd: "start-todo", IsLRO: true},
				&Command{MethodCmd: "list-todo"},
			},
		},
	}

	cmd := Command{
		Service:           name,
		ServiceClientType: name + "Client",
		MethodCmd:         strings.ToLower(name),
		ShortDesc:         "Sub-command for Service: " + name,
		Imports:           g.imports,
		EnvPrefix:         strings.ToUpper(g.root + "_" + name),
		SubCommands:       g.subcommands[name],
	}

	g.genServiceCmdFile(&cmd)
	if g.response.GetError() != "" {
		t.Errorf("Error generating the service_file: %s", g.response.GetError())
		return
	}

	file := g.response.File[0]

	if file.GetName() != "todo_service.go" {
		t.Errorf("(%+v).genServiceCmdFile(%+v) = %s, want %s", g, cmd, file.GetName(), "todo_service.go")
	}
	txtdiff.Diff(t, file.GetContent(), filepath.Join("testdata", "service_file.want"))
}
