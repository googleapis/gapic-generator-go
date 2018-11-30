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
	// CmdTemplate is the template for a cobra subcommand
	CmdTemplate = `{{$inputVar := (print .Method "Input")}}
{{$methodCmdVar := (print .Method "Cmd")}}
{{$pollingCmdVar := (print .Method "PollCmd")}}
{{$pollingOperationVar := (print .Method "PollOperation")}}
{{$fromFileVar := (print .Method "FromFile")}}
{{$outFileVar := (print .Method "OutFile")}}
{{$serviceCmdVar := (print .Service "ServiceCmd")}}
{{$followVar := (print .Method "Follow")}}
{{ $serviceClient := ( print .Service "Client" ) }}
// AUTO-GENERATED CODE. DO NOT EDIT.

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/golang/protobuf/jsonpb"
	{{ range $key, $pkg := .Imports}}
	{{ $pkg.Name }} "{{ $pkg.Path }}"
	{{ end }}
)

var {{ $inputVar }} {{ .InputMessage }}
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
var {{ ( $val.GenOneOfVarName $inputVar ) }} string
{{ range $oneOfKey, $oneOfVal := $val.OneOfs}}
var {{($oneOfVal.GenOneOfVarName $inputVar)}} {{if $oneOfVal.IsNested }}{{ $oneOfVal.MessageImport.Name }}.{{ $oneOfVal.Message }}{{ else }}{{ $.InputMessage }}{{ end }}_{{ ( title $oneOfKey ) }}
{{ end }}
{{ end }}
{{ range .Flags }}
{{ if and ( .IsMessage ) .Repeated }}
var {{ ( .GenOtherVarName $inputVar) }} []string
{{ else if ( .IsEnum ) }}
var {{ ( .GenOtherVarName $inputVar ) }} string
{{ end }}
{{ end }}

func init() {
	{{ $serviceCmdVar }}.AddCommand({{ $methodCmdVar }})
	{{ range .NestedMessages }}
	{{ $inputVar }}{{ .FieldName }} = new({{ .FieldType }})
	{{ end }}
	{{ range .Flags }}
	{{ $methodCmdVar }}.Flags().{{ (.GenFlag $inputVar) }}
	{{ end }}
	{{ range $key, $val := .OneOfSelectors }}
	{{ $methodCmdVar }}.Flags().{{ ($val.GenFlag $inputVar) }}
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
			{{ ( .GenRequired ) }}
			{{ end }}
			{{ end }}
			{{ range $key, $val := .OneOfSelectors }}
			{{ ( $val.GenRequired ) }}
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
			{{ if not .ClientStreaming }}
			err = jsonpb.Unmarshal(in, &{{ $inputVar }})
			if err != nil {
				return err
			}
			{{ end }}
		} {{ if or .OneOfSelectors .HasEnums }} else {
			{{ if .OneOfSelectors }}
			{{ range $key, $val := .OneOfSelectors }}
			switch {{ ( .GenOneOfVarName $inputVar ) }} {
			{{ range $oneOfKey, $oneOfVal := .OneOfs }}
			case "{{$oneOfKey}}":
				{{$inputVar}}.{{($val.InputFieldName)}} = &{{($oneOfVal.GenOneOfVarName $inputVar)}}
			{{ end }}
			default:
				return fmt.Errorf("Missing oneof choice for {{ .Name }}")
			}
			{{ end }}
			{{ end }}
			{{ if .HasEnums }}
			{{ range .Flags }}
			{{ if ( .IsEnum ) }}{{ $enumType := (print .MessageImport.Name "." .Message ) }}
			{{ $inputVar }}.{{ ( .InputFieldName ) }} = {{ $enumType }}({{ $enumType }}_value[strings.ToUpper({{ ( .GenOtherVarName $inputVar ) }})])
			{{ end }} 
			{{ end }}
			{{ end }}
		}
		{{ end }}
		{{ end }}
		{{ range .Flags }}
		{{ if and ( .IsMessage ) .Repeated }}
		// unmarshal JSON strings into slice of structs
		for _, item := range {{ ( .GenOtherVarName $inputVar) }} {
			tmp := {{ .MessageImport.Name }}.{{ .Message }}{}
			err = json.Unmarshal([]byte(item), &tmp)
			if err != nil {
				return
			}

			{{ if .IsOneOfField }}
			{{ $sliceAccessor := (print ( .GenOneOfVarName $inputVar) "." ( .OneOfInputFieldName )) }}
			{{ $sliceAccessor }} = append({{ $sliceAccessor }}, &tmp)
			{{ else }}
			{{ $sliceAccessor := (print $inputVar "." ( .InputFieldName )) }}
			{{ $sliceAccessor }} = append({{ $sliceAccessor }}, &tmp)
			{{ end }}
		}
		{{ end }}
		{{ end }}
		{{ if and (eq .OutputMessageType "") ( not .ClientStreaming ) }}
		if Verbose {
			printVerboseInput("{{ .Service }}", "{{ .Method }}", &{{ $inputVar }})
		}
		err = {{ $serviceClient }}.{{ .Method }}(ctx, &{{ $inputVar }})
		{{ else }}
		{{ if and ( not .ClientStreaming ) ( not .Paged ) }}
		if Verbose {
			printVerboseInput("{{ .Service }}", "{{ .Method }}", &{{ $inputVar }})
		}
		resp, err := {{ $serviceClient }}.{{ .Method }}(ctx, &{{ $inputVar }})
		{{ else if and .Paged ( not .IsLRO )}}
		if Verbose {
			printVerboseInput("{{ .Service }}", "{{ .Method }}", &{{ $inputVar }})
		}
		iter := {{ $serviceClient }}.{{ .Method }}(ctx, &{{ $inputVar }})
		{{ else if ( not .IsLRO )}}
		stream, err := {{ $serviceClient }}.{{ .Method }}(ctx)
		{{ else }}
		if Verbose {
			printVerboseInput("{{ .Service }}", "{{ .Method }}", &{{ $inputVar }})
		}
		resp, err := {{ $serviceClient }}.{{ .Method }}(ctx, &{{ $inputVar }})
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
					d, _ := json.MarshalIndent(res, "", "  ")
					str = string(d)
				}
				fmt.Fprintln(out, str)
			}
		}()
		{{ end }}
		if Verbose {
			fmt.Println("Client stream open. Close with blank line.")
		}

		scanner := bufio.NewScanner(in)
    for scanner.Scan() {
				input := scanner.Text()
				if input == "" {
					break
				}
        err = jsonpb.UnmarshalString(input, &{{ $inputVar }})
				if err != nil {
					return err
				}
				
				err = stream.Send(&{{ $inputVar }})
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
		page, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				fmt.Println("No more results")
				return nil
			}

			return err
		}

		data := make(map[string]interface{})
		data["page"] = page

		//get next page token
		_, err = iter.Next()
		if err != nil && err != iterator.Done {
			return err
		}
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

func (g *gcli) genCommands() {
	helpers := make(template.FuncMap)
	helpers["title"] = title

	t := template.Must(template.New("cmd").Funcs(helpers).Parse(CmdTemplate))

	for _, cmd := range g.commands {
		g.pt.Reset()

		t.Execute(g.pt.Writer(), cmd)

		g.addGoFile(cmd.MethodCmd + ".go")

		g.pt.Reset()
	}
}
