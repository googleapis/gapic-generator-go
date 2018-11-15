package gencli

import (
	"strings"
	"text/template"
)

const (
	// CmdTemplate is the template for a cobra subcommand
	CmdTemplate = `{{$inputVar := (print .Method "Input")}}
{{$methodCmdVar := (print .Method "Cmd")}}
{{$pollingCmdVar := (print .Method "PollCmd")}}
{{$pollingOperationVar := (print .Method "PollOperation")}}
{{$fromFileVar := (print .Method "FromFile")}}
{{$serviceCmdVar := (print .Service "ServiceCmd")}}
{{$followVar := (print .Method "Follow")}}
{{ $serviceClient := ( print .Service "Client" ) }}
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
{{ if .IsLRO }}
var {{ $followVar }} bool 

var {{ $pollingOperationVar }} string
{{ end }}
{{ range $key, $val := .OneOfSelectors }}
var {{ $inputVar }}{{ ($val.InputFieldName) }} string
{{ range $oneOfKey, $oneOfVal := $val.OneOfs}}
var {{($oneOfVal.GenOneOfVarName $inputVar)}} {{$.InputMessage}}_{{ ( title $oneOfKey ) }}
{{ end }}
{{ end }}
{{ range .Flags }}
{{ if and ( .IsMessage ) .Repeated }}
var {{ ( .GenRepeatedMessageVarName $inputVar) }} []string
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
		} {{ if .OneOfSelectors }} else {
			{{ range $key, $val := .OneOfSelectors }}
			switch {{ $inputVar }}{{ (.InputFieldName) }} {
			{{ range $oneOfKey, $oneOfVal := .OneOfs }}
			case "{{$oneOfKey}}":
				{{$inputVar}}.{{($val.InputFieldName)}} = &{{($oneOfVal.GenOneOfVarName $inputVar)}}
			{{ end }}
			default:
				return fmt.Errorf("Missing oneof choice for {{ .Name }}")
			}
			{{end}}
		}
		{{ end }}
		{{ end }}
		{{ range .Flags }}
		{{ if and ( .IsMessage ) .Repeated }}
		// unmarshal JSON strings into slice of structs
		for _, item := range {{ ( .GenRepeatedMessageVarName $inputVar) }} {
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
		{{ else if and .ClientStreaming ( not .ServerStreaming ) }}
		if Verbose {
			fmt.Println("Client stream open. Close with blank line.")
		}

		scanner := bufio.NewScanner(in)
    for err == nil && scanner.Scan() {
				input := scanner.Text()
				if input == "" {
					break
				}
        err = jsonpb.UnmarshalString(input, &{{ $inputVar }})
				if err != nil {
					return err
				}
				
				err = stream.Send(&{{ $inputVar }})
    }
    if err := scanner.Err(); err != nil {
        return err
    }
		
		resp, err := stream.CloseAndRecv()
		if err != nil {
			return err
		}
		{{ if not .IsLRO }}
		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(resp){{end}}
		{{ else if and .ClientStreaming .ServerStreaming }}
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
	helpers["title"] = strings.Title

	t := template.Must(template.New("cmd").Funcs(helpers).Parse(CmdTemplate))

	for _, cmd := range g.commands {
		g.pt.Reset()

		t.Execute(g.pt.Writer(), cmd)

		g.addGoFile(cmd.MethodCmd + ".go")

		g.pt.Reset()
	}
}
