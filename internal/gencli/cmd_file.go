package gencli

import (
	"text/template"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

const (
	// CmdTemplate is the template for a cobra subcommand
	CmdTemplate = `package main

import (
	"fmt"

	"github.com/spf13/cobra"
)
{{$methodCmdVar := (print .Method "Cmd")}}
{{$fromFileVar := (print .Method "FromFile")}}
{{ range .Flags }}
{{ (.GenFlagVar) }}
{{ end }}
{{ if .Flags }}
var {{ $fromFileVar }} string
{{ end }}

func init() {
	rootCmd.AddCommand({{ $methodCmdVar }})
	{{ range .Flags }}
	{{ $methodCmdVar }}.Flags().{{ (.GenFlag ) }}
	{{ end }}
	{{ if .Flags }}
	{{ $methodCmdVar }}.Flags().StringVar(&{{ $fromFileVar }}, "from_file", "", "")
	{{ end }}
}

var {{$methodCmdVar}} = &cobra.Command{
  Use:   "{{ .MethodCmd }}",
  {{ if (ne .ShortDesc "") }}Short: "{{ .ShortDesc }}",{{ end }}
	{{ if (ne .LongDesc "") }}Long: {{ .LongDesc }},{{ end }}
	PreRun: func(cmd *cobra.Command, args []string) {
		{{ if .Flags }}
		if {{ $fromFileVar }} == "" {
			{{ range .Flags }}
			{{ if .Required }}
			{{ ( .GenRequired ) }}
			{{ end }}
			{{ end }}
		}
		{{ end }}
	},
	Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Hello, from {{ .Method }}")
  },
}
	`
)

func (g *gcli) genCommands() {
	for _, cmd := range g.commands {
		g.pt.Reset()
		template.Must(template.New("cmd").Parse(CmdTemplate)).Execute(g.pt.Writer(), cmd)

		g.response.File = append(g.response.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(cmd.MethodCmd + ".go"),
			Content: proto.String(g.pt.String()),
		})

		g.pt.Reset()
	}
}
