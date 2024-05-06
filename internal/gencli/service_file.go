// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gencli

import (
	"text/template"
)

const (

	// serviceTemplate is the template string for generated {service}.go
	serviceTemplate = `{{ $serviceCmdVar := (print .Service "ServiceCmd") }}
{{ $serviceClient := ( print .Service "Client" ) }}
{{ $serviceSubCommands := (print .Service "SubCommands" ) }}
{{ $serviceConfig := (print .Service "Config" ) }}
package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"golang.org/x/oauth2"
	{{ range $key, $pkg := .Imports}}
	{{ $pkg.Name }} "{{ $pkg.Path }}"
	{{ end }}
)

var {{ $serviceConfig }} *viper.Viper
var {{ $serviceClient }} *gapic.{{.ServiceClientType}}
var {{ $serviceSubCommands }} []string = []string{
	{{ range .SubCommands }}"{{ .MethodCmd }}",
	{{ if .IsLRO }}"poll-{{ .MethodCmd }}",{{ end }}{{ end }}
}

func init() {
	rootCmd.AddCommand({{ $serviceCmdVar }})

	{{ $serviceConfig }} = viper.New()
	{{ $serviceConfig }}.SetEnvPrefix("{{ .EnvPrefix }}")
	{{ $serviceConfig }}.AutomaticEnv()

	{{ $serviceCmdVar }}.PersistentFlags().Bool("insecure", false, "Make insecure client connection. Or use {{.EnvPrefix}}_INSECURE. Must be used with \"address\" option")
	{{ $serviceConfig }}.BindPFlag("insecure", {{ $serviceCmdVar }}.PersistentFlags().Lookup("insecure"))
	{{ $serviceConfig }}.BindEnv("insecure")

	{{ $serviceCmdVar }}.PersistentFlags().String("address", "", "Set API address used by client. Or use {{.EnvPrefix}}_ADDRESS.")
	{{ $serviceConfig }}.BindPFlag("address", {{ $serviceCmdVar }}.PersistentFlags().Lookup("address"))
	{{ $serviceConfig }}.BindEnv("address")

	{{ $serviceCmdVar }}.PersistentFlags().String("token", "", "Set Bearer token used by the client. Or use {{.EnvPrefix}}_TOKEN.")
	{{ $serviceConfig }}.BindPFlag("token", {{ $serviceCmdVar }}.PersistentFlags().Lookup("token"))
	{{ $serviceConfig }}.BindEnv("token")

	{{ $serviceCmdVar }}.PersistentFlags().String("api_key", "", "Set API Key used by the client. Or use {{.EnvPrefix}}_API_KEY.")
	{{ $serviceConfig }}.BindPFlag("api_key", {{ $serviceCmdVar }}.PersistentFlags().Lookup("api_key"))
	{{ $serviceConfig }}.BindEnv("api_key")
}

var {{ $serviceCmdVar }} = &cobra.Command{
	Use:   "{{ .MethodCmd }}",
	{{ if (ne .ShortDesc "") }}Short: "{{ .ShortDesc }}",{{ end }}
	{{ if (ne .LongDesc "") }}Long: "{{ .LongDesc }}",{{ end }}
	ValidArgs: {{ $serviceSubCommands }},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		var opts []option.ClientOption

		address := {{ $serviceConfig }}.GetString("address")
		if address != "" {
			opts = append(opts, option.WithEndpoint(address))
		}

		if {{ $serviceConfig }}.GetBool("insecure"){
			if address == "" {
				return fmt.Errorf("Missing address to use with insecure connection")
			}

			conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return err
			}
			opts = append(opts, option.WithGRPCConn(conn))
		}

		if token := {{ $serviceConfig }}.GetString("token"); token != "" {
			opts = append(opts, option.WithTokenSource(oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: token,
					TokenType:   "Bearer",
				})))
		}

		if key := {{ $serviceConfig }}.GetString("api_key"); key != "" {
			opts = append(opts, option.WithAPIKey(key))
		}

		{{ $serviceClient }}, err = gapic.New{{.ServiceClientType}}(ctx, opts...)
		return
	},
}
`
)

var serviceTemplateCompiled *template.Template

func init() {
	serviceTemplateCompiled = template.Must(template.New("service").Parse(serviceTemplate))
}

func (g *gcli) genServiceCmdFile(cmd *Command) {

	g.pt.Reset()

	g.pt.Printf("// Code generated. DO NOT EDIT.\n")

	serviceTemplateCompiled.Execute(g.pt.Writer(), cmd)

	g.addGoFile(cmd.MethodCmd + "_service" + ".go")
}
