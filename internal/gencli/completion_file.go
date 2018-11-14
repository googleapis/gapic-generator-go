package gencli

import (
	"strings"
	"text/template"
)

const (
	// CompletionTemplate is the template string for the bash completion generation command
	CompletionTemplate = `package main

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
	template.Must(template.New("comp").Parse(CompletionTemplate)).Execute(g.pt.Writer(), Command{
		MethodCmd: strings.ToLower(g.root),
	})

	g.addGoFile("completion.go")

	g.pt.Reset()
}
