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
	// completionTemplate is the template string for the bash completion generation command
	completionTemplate = `
	package main

	import (
		"os"
	
		"github.com/spf13/cobra"
	)
	
	func init() {
		rootCmd.AddCommand(completionCmd)
	}
	
	// completionCmd represents the completion command
	var completionCmd = &cobra.Command{
		Use:   "completion",
		Short: "Emits bash a completion for {{ .MethodCmd }}",
		Long: {{ .BackTick }}Enable bash completion like so:
		Linux:
			source <({{ .MethodCmd }} completion)
		Mac:
			brew install bash-completion
			{{ .MethodCmd }} completion > $(brew --prefix)/etc/bash_completion.d/{{ .MethodCmd }}{{ .BackTick }},
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletion(os.Stdout)
		},
	}`
)

func (g *gcli) genCompletionCmdFile() {
	g.pt.Reset()

	g.pt.Printf("// Code generated. DO NOT EDIT.\n")
	template.Must(template.New("comp").Parse(completionTemplate)).Execute(g.pt.Writer(), struct {
		MethodCmd string
		BackTick  string
	}{
		MethodCmd: strings.ToLower(g.root),
		BackTick:  "`",
	})

	g.addGoFile("completion.go")
}
