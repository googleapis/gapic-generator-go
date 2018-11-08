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

	"golang.org/x/net/context"
	"github.com/spf13/cobra"
	{{ range $key, $pkg := .Imports}}
	{{ $pkg.Name }} "{{ $pkg.Path }}"
	{{ end }}
)
{{ $serviceCmdVar := (print .Service "Cmd") }}

var client *gapic.{{.Service}}Client
var ctx context.Context

func init() {
	rootCmd.AddCommand({{ $serviceCmdVar }})
}

var {{ $serviceCmdVar }} = &cobra.Command{
	Use:   "{{ .MethodCmd }}",
	{{ if (ne .ShortDesc "") }}Short: "{{ .ShortDesc }}",{{ end }}
	{{ if (ne .LongDesc "") }}Long: {{ .LongDesc }},{{ end }}
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx = context.Background()
		client, err = gapic.New{{.Service}}Client(ctx)
		return
	},
}
`
)

func (g *gcli) genServiceCmdFiles() {
	g.pt.Reset()
	t := template.Must(template.New("service").Parse(ServiceTemplate))

	for _, srv := range g.services {
		name := pbinfo.ReduceServName(srv.GetName(), "")
		cmd := Command{
			Service:   name,
			MethodCmd: strings.ToLower(name),
			ShortDesc: "Sub-command for Service: " + name,
			Imports:   g.imports,
		}

		// add any available comment as usage
		key := pbinfo.BuildElementCommentKey(g.descInfo.ParentFile[srv], srv)
		if cmt, ok := g.descInfo.Comments[key]; ok {
			cmt = strings.TrimSpace(strings.Replace(cmt, "\n", " ", -1))

			cmd.LongDesc = cmt
			cmd.ShortDesc = toShortUsage(cmt)
		}

		t.Execute(g.pt.Writer(), cmd)

		g.addGoFile(cmd.MethodCmd + ".go")

		g.pt.Reset()
	}
}
