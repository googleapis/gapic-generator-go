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

	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var Verbose, OutputJSON bool
var ctx context.Context

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Print verbose output")
	rootCmd.PersistentFlags().BoolVarP(&OutputJSON, "json", "j", false, "Print JSON output")
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

func printVerboseInput(srv, mthd string, data interface{}) {
	fmt.Println("Service:", srv)
	fmt.Println("Method:", mthd)
	fmt.Print("Input: ")
	printMessage(data)
}

func printMessage(data interface{}) {
	var s string

	if OutputJSON {
		d, _ := json.MarshalIndent(data, "", "  ")
		s = string(d)
	} else if msg, ok := data.(proto.Message); ok {
		s = msg.String()
	} else if page, ok := data.(map[string]interface{}); ok {
		s = fmt.Sprintf("%v", page)
	}

	fmt.Println(s)
}
`
)

func (g *gcli) genRootCmdFile() {
	g.pt.Reset()
	name := strings.ToLower(g.root)
	template.Must(template.New("root").Parse(RootTemplate)).Execute(g.pt.Writer(), Command{
		MethodCmd: name,
		ShortDesc: "Root command of " + g.root,
	})

	g.addGoFile(name + ".go")

	g.pt.Reset()
}
