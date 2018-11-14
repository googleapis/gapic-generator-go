package gencli

import (
	"strings"
	"text/template"

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

const (

	// ServiceTemplate is the template string for generated {service}.go
	ServiceTemplate = `{{ $serviceCmdVar := (print .Service "ServiceCmd") }}
{{ $serviceClient := ( print .Service "Client" ) }}
{{ $serviceSubCommands := (print .Service "SubCommands" )}}
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	{{ range $key, $pkg := .Imports}}
	{{ $pkg.Name }} "{{ $pkg.Path }}"
	{{ end }}
)

var {{ $serviceClient }} *gapic.{{.Service}}Client
var {{ $serviceSubCommands }} []string = []string{
	{{ range .SubCommands }}"{{ .MethodCmd }}",
	{{ if .IsLRO }}"poll-{{ .MethodCmd }}",{{ end }}{{ end }}
}

func init() {
	rootCmd.AddCommand({{ $serviceCmdVar }})
}

var {{ $serviceCmdVar }} = &cobra.Command{
	Use:   "{{ .MethodCmd }}",
	{{ if (ne .ShortDesc "") }}Short: "{{ .ShortDesc }}",{{ end }}
	{{ if (ne .LongDesc "") }}Long: "{{ .LongDesc }}",{{ end }}
	ValidArgs: {{ $serviceSubCommands }},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		var opts []option.ClientOption

		if address := os.Getenv("{{.EnvPrefix}}_ADDRESS"); address != "" {
			opts = append(opts, option.WithEndpoint(address))

			if Insecure {
				conn, err := grpc.Dial(address, grpc.WithInsecure())
				if err != nil {
					return err
				}
				opts = append(opts, option.WithGRPCConn(conn))
			}
		}

		ctx = context.Background()
		{{ $serviceClient }}, err = gapic.New{{.Service}}Client(ctx, opts...)
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
			Service:     name,
			MethodCmd:   strings.ToLower(name),
			ShortDesc:   "Sub-command for Service: " + name,
			Imports:     g.imports,
			EnvPrefix:   strings.ToUpper(g.root + "_" + name),
			SubCommands: g.subcommands[srv.GetName()],
		}

		// add any available comment as usage
		key := pbinfo.BuildElementCommentKey(g.descInfo.ParentFile[srv], srv)
		if cmt, ok := g.descInfo.Comments[key]; ok {
			cmt = sanitizeComment(cmt)

			cmd.LongDesc = cmt
			cmd.ShortDesc = toShortUsage(cmt)
		}

		t.Execute(g.pt.Writer(), cmd)

		g.addGoFile(cmd.MethodCmd + "_service" + ".go")

		g.pt.Reset()
	}
}
