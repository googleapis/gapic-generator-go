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
	"strings"
	"text/template"
)

const (

	// rootTemplate is the template string for generated root.go
	rootTemplate = `
package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
)

var Verbose, OutputJSON bool
var ctx = context.Background()
var marshaler = &jsonpb.Marshaler{Indent: "  "}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Print verbose output")
	rootCmd.PersistentFlags().BoolVarP(&OutputJSON, "json", "j", false, "Print JSON output")
}

var rootCmd = &cobra.Command{
	Use:   "{{ .MethodCmd }}",
	{{ if (ne .ShortDesc "") }}Short: "{{ .ShortDesc }}",{{ end }}
	{{ if (ne .LongDesc "") }}Long: "{{ .LongDesc }}",{{ end }}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}

func printVerboseInput(srv, mthd string, data interface{}) {
	fmt.Println("Service:", srv)
	fmt.Println("Method:", mthd)
	fmt.Print("Input: ")
	printMessage(data)
}

func printMessage(data interface{}) {
	var s string

	if msg, ok := data.(proto.Message); ok {
		s = msg.String()
		if OutputJSON {
			var b bytes.Buffer
			marshaler.Marshal(&b, msg)
			s = b.String()
		}
	}

	fmt.Println(s)
}
`
)

func (g *gcli) genRootCmdFile() {
	g.pt.Reset()

	g.pt.Printf("// Code generated. DO NOT EDIT.\n")

	name := strings.ToLower(g.root)
	template.Must(template.New("root").Parse(rootTemplate)).Execute(g.pt.Writer(), Command{
		MethodCmd: name,
		ShortDesc: "Root command of " + g.root,
	})

	g.addGoFile(name + ".go")
}
