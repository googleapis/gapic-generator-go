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
	completionTemplate = `// AUTO-GENERATED CODE. DO NOT EDIT.
	
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
		Long: ` + "`Enable bash completion like so:\n" +
		"Linux:\n" +
		"  source <({{ .MethodCmd }} completion)\n" +
		"Mac:\n" +
		"  brew install bash-completion\n" +
		"  {{ .MethodCmd }} completion > $(brew --prefix)/etc/bash_completion.d/{{ .MethodCmd }}`,\n" +
		`Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletion(os.Stdout)
		},
	}`
)

func (g *gcli) genCompletionCmdFile() {
	g.pt.Reset()
	template.Must(template.New("comp").Parse(completionTemplate)).Execute(g.pt.Writer(), Command{
		MethodCmd: strings.ToLower(g.root),
	})

	g.addGoFile("completion.go")
}
