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
	// cmdTemplate is the template for a cobra subcommand
	cmdTemplate = `{{$methodCmdVar := (print .Method "Cmd")}}
{{$pollingCmdVar := (print .Method "PollCmd")}}
{{$pollingOperationVar := (print .Method "PollOperation")}}
{{$fromFileVar := (print .Method "FromFile")}}
{{$outFileVar := (print .Method "OutFile")}}
{{$serviceCmdVar := (print .Service "ServiceCmd")}}
{{$followVar := (print .Method "Follow")}}
{{ $serviceClient := ( print .Service "Client" ) }}
package main

import (
	"github.com/spf13/cobra"
	{{ range $key, $pkg := .Imports}}
	{{ $pkg.Name }} "{{ $pkg.Path }}"
	{{ end }}
)
{{ if not .ClientStreaming }}
var {{ .InputMessageVar }} {{ .InputMessage }}{{ end }}
{{ if or .Flags .ClientStreaming }}
var {{ $fromFileVar }} string
{{ end }}
{{ if and .ServerStreaming .ClientStreaming }}
var {{ $outFileVar }} string
{{ end }}
{{ if .IsLRO }}
var {{ $followVar }} bool 

var {{ $pollingOperationVar }} string
{{ end }}
{{ range $key, $val := .OneOfSelectors }}
var {{ $val.VarName }} string
{{ range $oneOfKey, $oneOfVal := $val.OneOfs}}
var {{$oneOfVal.VarName}} {{if $oneOfVal.IsNested }}{{ $oneOfVal.MessageImport.Name }}.{{ $oneOfVal.Message }}{{ else }}{{ $.InputMessage }}{{ end }}_{{ ( title $oneOfKey ) }}
{{ end }}
{{ end }}
{{ range .Flags }}
{{ if and ( .IsMessage ) .Repeated }}
var {{ .VarName }} []string
{{ else if ( .IsEnum ) }}
var {{ .VarName }} {{if .Repeated }}[]{{ end }}string
{{ end }}
{{ end }}

func init() {
	{{ $serviceCmdVar }}.AddCommand({{ $methodCmdVar }})
	{{ range .NestedMessages }}
	{{ .FieldName }} = new({{ .FieldType }})
	{{ end }}
	{{ range .Flags }}
	{{ $methodCmdVar }}.Flags().{{ (.GenFlag) }}
	{{ end }}
	{{ range $key, $val := .OneOfSelectors }}
	{{ $methodCmdVar }}.Flags().{{ ($val.GenFlag) }}
	{{ end }}
	{{ if or .Flags .ClientStreaming }}
	{{ $methodCmdVar }}.Flags().StringVar(&{{ $fromFileVar }}, "from_file", "", "Absolute path to JSON file containing request payload")
	{{ end }}
	{{ if and .ClientStreaming .ServerStreaming }}
	{{ $methodCmdVar }}.Flags().StringVar(&{{ $outFileVar }}, "out_file", "", "Absolute path to a file to pipe output to")
	{{ $methodCmdVar }}.MarkFlagRequired("out_file")
	{{ end }}
	{{ if .IsLRO }}
	{{ $methodCmdVar }}.Flags().BoolVar(&{{ $followVar }}, "follow", false, "Block until the long running operation completes")

	{{ $serviceCmdVar }}.AddCommand({{ $pollingCmdVar }})

	{{ $pollingCmdVar }}.Flags().BoolVar(&{{ $followVar }}, "follow", false, "Block until the long running operation completes")

	{{ $pollingCmdVar }}.Flags().StringVar(&{{$pollingOperationVar}}, "operation", "", "Required. Operation name to poll for")

	{{ $pollingCmdVar }}.MarkFlagRequired("operation")

	{{ end }}
}

var {{$methodCmdVar}} = &cobra.Command{
  Use:   "{{ .MethodCmd }}",
  {{ if (ne .ShortDesc "") }}Short: "{{ .ShortDesc }}",{{ end }}
	{{ if (ne .LongDesc "") }}Long: "{{ .LongDesc }}",{{ end }}
	PreRun: func(cmd *cobra.Command, args []string) {
		{{ if or .Flags .OneOfSelectors }}
		if {{ $fromFileVar }} == "" {
			{{ range .Flags }}
			{{ if and .Required ( not .IsOneOfField ) }}
			cmd.MarkFlagRequired("{{ .Name }}")
			{{ end }}
			{{ end }}
			{{ range $key, $val := .OneOfSelectors }}
			cmd.MarkFlagRequired("{{ $val.Name }}")
			{{ end }}
		}
		{{ end }}
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		{{ if or .Flags .ClientStreaming }}
		in := os.Stdin
		if {{ $fromFileVar }} != "" {
			in, err = os.Open({{ $fromFileVar }})
			if err != nil {
				return err
			}
			defer in.Close()
			{{ if not .ClientStreaming }}
			err = jsonpb.Unmarshal(in, &{{ .InputMessageVar }})
			if err != nil {
				return err
			}
			{{ end }}
		} {{ if or .OneOfSelectors .HasEnums }} else {
			{{ if .OneOfSelectors }}
			{{ range $key, $val := .OneOfSelectors }}
			switch {{ .VarName }} {
			{{ range $oneOfKey, $oneOfVal := .OneOfs }}
			case "{{$oneOfKey}}":
				{{ $.InputMessageVar }}.{{$val.FieldName}} = &{{$oneOfVal.VarName}}
			{{ end }}
			default:
				return fmt.Errorf("Missing oneof choice for {{ .Name }}")
			}
			{{ end }}
			{{ end }}
			{{ if .HasEnums }}
			{{ range .Flags }}
			{{ if ( .IsEnum ) }}{{ $enumType := (print .MessageImport.Name "." .Message ) }}{{ $requestField := ( .EnumFieldAccess $.InputMessageVar ) }}
			{{ if .Repeated }}
			for _, in := range {{ .VarName }} {
				val := {{ $enumType }}({{ $enumType }}_value[strings.ToUpper(in)])
				{{ $requestField }} = append({{ $requestField }}, val)
			}
			{{ else }}
			{{ $requestField }} = {{ $enumType }}({{ $enumType }}_value[strings.ToUpper({{ .VarName }})])
			{{ end }}
			{{ end }} 
			{{ end }}
			{{ end }}
		}
		{{ end }}
		{{ end }}
		{{ range .Flags }}
		{{ if and ( .IsMessage ) .Repeated }}
		// unmarshal JSON strings into slice of structs
		for _, item := range {{ .VarName }} {
			tmp := {{ .MessageImport.Name }}.{{ .Message }}{}
			err = jsonpb.UnmarshalString(item, &tmp)
			if err != nil {
				return
			}

			{{ .SliceAccessor }} = append({{ .SliceAccessor }}, &tmp)
		}
		{{ end }}
		{{ end }}
		{{ if and (eq .OutputMessageType "") ( not .ClientStreaming ) }}
		if Verbose {
			printVerboseInput("{{ .Service }}", "{{ .Method }}", &{{ .InputMessageVar }})
		}
		err = {{ $serviceClient }}.{{ .Method }}(ctx, &{{ .InputMessageVar }})
		{{ else }}
		{{ if and ( not .ClientStreaming ) ( not .Paged ) }}
		if Verbose {
			printVerboseInput("{{ .Service }}", "{{ .Method }}", &{{ .InputMessageVar }})
		}
		resp, err := {{ $serviceClient }}.{{ .Method }}(ctx, &{{ .InputMessageVar }})
		{{ else if and .Paged ( not .IsLRO )}}
		if Verbose {
			printVerboseInput("{{ .Service }}", "{{ .Method }}", &{{ .InputMessageVar }})
		}
		iter := {{ $serviceClient }}.{{ .Method }}(ctx, &{{ .InputMessageVar }})
		{{ else if ( not .IsLRO )}}
		stream, err := {{ $serviceClient }}.{{ .Method }}(ctx)
		{{ else }}
		if Verbose {
			printVerboseInput("{{ .Service }}", "{{ .Method }}", &{{ .InputMessageVar }})
		}
		resp, err := {{ $serviceClient }}.{{ .Method }}(ctx, &{{ .InputMessageVar }})
		{{ end }}
		{{ if and .ServerStreaming ( not .ClientStreaming ) }}
		var item *{{ .OutputMessageType }}
		for {
			item, err = resp.Recv()
			if err != nil {
				break
			}

			if Verbose {
				fmt.Print("Output: ")
			}
			printMessage(item)
		}

		if err == io.EOF {
			return nil
		}
		{{ else if .ClientStreaming }}
		{{ if .ServerStreaming }}
		out, err := os.OpenFile({{ $outFileVar}}, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return err
		}

		// start background stream receive
		go func() {
			var res *{{ .OutputMessageType }}
			for {
				res, err = stream.Recv()
				if err != nil {
					return
				}

				str := res.String()
				if OutputJSON {
					str, _ = marshaler.MarshalToString(res)
				}
				fmt.Fprintln(out, str)
			}
		}()
		{{ end }}
		if Verbose {
			fmt.Println("Client stream open. Close with ctrl+D.")
		}

		var {{.InputMessageVar }} {{ .InputMessage }}
		scanner := bufio.NewScanner(in)
    for scanner.Scan() {
				input := scanner.Text()
				if input == "" {
					continue
				}
        err = jsonpb.UnmarshalString(input, &{{ .InputMessageVar }})
				if err != nil {
					return err
				}
				
				err = stream.Send(&{{ .InputMessageVar }})
				if err != nil {
					return err
				}
    }
    if err = scanner.Err(); err != nil {
        return err
    }
		{{ if .ServerStreaming }}
		err = stream.CloseSend()
		{{ else }}
		resp, err := stream.CloseAndRecv()
		if err != nil {
			return err
		}
		{{ if not .IsLRO }}
		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)
		{{ end }}
		{{ end }}
		{{ else if and .Paged (not .IsLRO ) }}
		// get requested page
		var items []interface{}
		data := make(map[string]interface{})

		// PageSize could be an integer with a specific precision.
		// Doing standard i := 0; i < PageSize; i++ creates i as
		// an int, creating a potential type mismatch. 
		for i := {{ .InputMessageVar }}.PageSize; i > 0; i-- {
			item, err := iter.Next()
			if err == iterator.Done {
				err = nil
				break
			} else if err != nil {
				return err
			}

			items = append(items, item)
		}

		data["page"] = items
		data["nextToken"] = iter.PageInfo().Token

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(data)
		{{ else if not .IsLRO }}
		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp)
		{{ end }}

		{{ if .IsLRO }}
		if !{{ $followVar }} {
			var s interface{}
			s = resp.Name()

			if OutputJSON {
				d := make(map[string]string)
				d["operation"] = resp.Name()
				s = d
			}

			printMessage(s)
			return err
		}

		result, err := resp.Wait(ctx)
		if err != nil {
			return err
		}

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(result)
		{{ end }}
		{{ end }}
		return err
  },
}

{{ if .IsLRO }}
var {{ $pollingCmdVar }} = &cobra.Command{
	Use: "poll-{{ .MethodCmd }}",
	Short: "Poll the status of a {{ .Method }}Operation by name",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		op := {{ $serviceClient }}.{{ .Method }}Operation({{ $pollingOperationVar }})

		if {{ $followVar }} {
			resp, err := op.Wait(ctx)
			if err != nil {
				return err
			}

			if Verbose {
				fmt.Print("Output: ")
			}
			printMessage(resp)
			return err
		}

		resp, err := op.Poll(ctx)
		if err != nil {
			return err
		} else if resp != nil {
			if Verbose {
				fmt.Print("Output: ")
			}

			printMessage(resp)
			return
		}

		fmt.Println(fmt.Sprintf("Operation %s not done", op.Name()))

		return err
	},
}
{{ end }}
`
)

var cmdTemplateCompiled *template.Template

func init() {
	helpers := make(template.FuncMap)
	helpers["title"] = title

	cmdTemplateCompiled = template.Must(template.New("cmd").Funcs(helpers).Parse(cmdTemplate))
}

func (g *gcli) genCommandFile(cmd *Command) {
	g.pt.Reset()

	g.pt.Printf("// Code generated. DO NOT EDIT.\n")

	cmdTemplateCompiled.Execute(g.pt.Writer(), cmd)

	g.addGoFile(cmd.MethodCmd + ".go")
}
