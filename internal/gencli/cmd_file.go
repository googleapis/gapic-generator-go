package gencli

import (
	"text/template"
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
{{$serviceCmdVar := (print .Service "Cmd")}}
{{ range .Flags }}
{{ (.GenFlagVar) }}
{{ end }}
{{ if .Flags }}
var {{ $fromFileVar }} string
{{ end }}

func init() {
	{{ $serviceCmdVar }}.AddCommand({{ $methodCmdVar }})
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
	{{ if (ne .LongDesc "") }}Long: "{{ .LongDesc }}",{{ end }}
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

		g.addGoFile(cmd.MethodCmd + ".go")

		g.pt.Reset()
	}
}
