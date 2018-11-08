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

var Verbose, Insecure bool

func init() {
	rootCmd.PersistentFlags().BoolVar(&Insecure, "insecure", false, "Make insecure client connection. Must be used with {{.EnvPrefix}}_ADDRESS.")
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

func main() {
	Execute()
}
`
)

func (g *gcli) genRootCmdFile() {
	g.pt.Reset()
	template.Must(template.New("root").Parse(RootTemplate)).Execute(g.pt.Writer(), Command{
		MethodCmd: strings.ToLower(g.Root),
		ShortDesc: "Root command of " + g.Root,
		EnvPrefix: strings.ToUpper(g.Root + "_{SERVICE}"),
	})

	g.addGoFile("root.go")

	g.pt.Reset()
}
