package gencli

import (
	"strings"
	"text/template"
)

const (

	// RootTemplate is the template string for generated root.go
	RootTemplate = `package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Verbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Print verbose output")
}

var rootCmd = &cobra.Command{
	Use:   "{{ .MethodCmd }}",
	{{ if (ne .ShortDesc "") }}Short: "{{ .ShortDesc }}",{{ end }}
	{{ if (ne .LongDesc "") }}Long: {{ .LongDesc }},{{ end }}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
`
)

func (g *gcli) genRootCmdFile(root string) {
	g.pt.Reset()
	template.Must(template.New("root").Parse(RootTemplate)).Execute(g.pt.Writer(), Command{
		MethodCmd: strings.ToLower(root),
		ShortDesc: "Root command of " + root,
	})

	g.addGoFile("root.go")

	g.pt.Reset()
}
