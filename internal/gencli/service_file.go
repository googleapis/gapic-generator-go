package gencli

import (
	"strings"
	"text/template"

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

const (

	// ServiceTemplate is the template string for generated {service}.go
	ServiceTemplate = `package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)
{{ $serviceCmdVar := (print .Service "Cmd") }}

func init() {
	rootCmd.AddCommand({{ $serviceCmdVar }})
}

var {{ $serviceCmdVar }} = &cobra.Command{
	Use:   "{{ .MethodCmd }}",
	{{ if (ne .ShortDesc "") }}Short: "{{ .ShortDesc }}",{{ end }}
	{{ if (ne .LongDesc "") }}Long: {{ .LongDesc }},{{ end }}
}
`
)

func (g *gcli) genServiceCmdFiles() {
	g.pt.Reset()
	t := template.Must(template.New("service").Parse(ServiceTemplate))

	for _, srv := range g.services {
		name := pbinfo.ReduceServName(srv.GetName(), "")
		lower := strings.ToLower(name)

		t.Execute(g.pt.Writer(), Command{
			Service:   name,
			MethodCmd: lower,
			ShortDesc: "Top level command for Service: " + name,
		})

		g.addGoFile(lower + ".go")

		g.pt.Reset()
	}
}
