package gencli

import (
	"text/template"
)

const (
	// CmdTemplate is the template for a cobra subcommand
	CmdTemplate = `{{$inputVar := (print .Method "Input")}}
{{$methodCmdVar := (print .Method "Cmd")}}
{{$fromFileVar := (print .Method "FromFile")}}
{{$serviceCmdVar := (print .Service "Cmd")}}
package main

import (
	"fmt"
	"encoding/json"

	"github.com/spf13/cobra"
	{{ range $key, $pkg := .Imports}}
	{{ $pkg.Name }} "{{ $pkg.Path }}"
	{{ end }}
)

var {{ $inputVar }} {{ .InputMessage }}
{{ if .Flags }}
var {{ $fromFileVar }} string
{{ range .Flags }}
{{ if and ( .IsMessage ) .Repeated }}
{{ ( .GenRepeatedMessageFlagVar $inputVar) }}
{{ end }}
{{ end }}
{{ end }}

func init() {
	{{ $serviceCmdVar }}.AddCommand({{ $methodCmdVar }})
	{{ range .NestedMessages }}
	{{ $inputVar }}.{{ .FieldName }} = new({{ .FieldType }})
	{{ end }}
	{{ range .Flags }}
	{{ $methodCmdVar }}.Flags().{{ (.GenFlag $inputVar) }}
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		{{ range .Flags }}
		{{ if and ( .IsMessage ) .Repeated }}
		{{ $sliceAccessor := (print $inputVar "." ( .InputFieldName )) }}
		// unmarshal JSON strings into slice of structs
		for _, item := range {{ $inputVar }}{{ ( .InputFieldName ) }} {
			tmp := {{ .MessageImport.Name }}.{{ .Message }}{}
			err = json.Unmarshal([]byte(item), &tmp)
			if err != nil {
				return
			}

			{{ $sliceAccessor }} = append({{ $sliceAccessor }}, &tmp)
		}
		{{ end }}
		{{ end }}
		fmt.Println("Hello, from {{ .Method }}")
		{{ if (eq .OutputType "") }}
		err = client.{{ .Method }}(ctx, &{{ $inputVar }})
		{{ else }}
		resp, err := client.{{ .Method }}(ctx, &{{ $inputVar }})
		fmt.Println(resp)
		{{ end }}
		return err
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
