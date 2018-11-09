package gencli

import (
	"text/template"
)

const (
	// CmdTemplate is the template for a cobra subcommand
	CmdTemplate = `{{$inputVar := (print .Method "Input")}}
{{$methodCmdVar := (print .Method "Cmd")}}
{{$fromFileVar := (print .Method "FromFile")}}
{{$serviceCmdVar := (print .Service "ServiceCmd")}}
package main

import (
	"io/ioutil"
	"encoding/json"

	"github.com/spf13/cobra"
	{{ range $key, $pkg := .Imports}}
	{{ $pkg.Name }} "{{ $pkg.Path }}"
	{{ end }}
)

var {{ $inputVar }} {{ .InputMessage }}
{{ if and .Flags ( not .ClientStreaming ) }}
var {{ $fromFileVar }} string
{{ end }}
{{ range .Flags }}
{{ if and ( .IsMessage ) .Repeated }}
var {{ ( .GenRepeatedMessageVarName $inputVar) }} []string
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
	{{ if and .Flags ( not .ClientStreaming ) }}
	{{ $methodCmdVar }}.Flags().StringVar(&{{ $fromFileVar }}, "from_file", "", "Absolute path to JSON file containing request payload")
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
		for _, item := range {{ ( .GenRepeatedMessageVarName $inputVar) }} {
			tmp := {{ .MessageImport.Name }}.{{ .Message }}{}
			err = json.Unmarshal([]byte(item), &tmp)
			if err != nil {
				return
			}

			{{ $sliceAccessor }} = append({{ $sliceAccessor }}, &tmp)
		}
		{{ end }}
		{{ end }}
		{{ if and .Flags ( not .ClientStreaming ) }}
		if {{ $fromFileVar }} != "" {
			data, err := ioutil.ReadFile({{ $fromFileVar }})
			if err != nil {
				return err
			}

			err = json.Unmarshal(data, &{{ $inputVar }})
			if err != nil {
				return err
			}
		}
		{{ end }}
		{{ if and (eq .OutputType "") ( not .ClientStreaming ) }}
		err = client.{{ .Method }}(ctx, &{{ $inputVar }})
		{{ else }}
		{{ if and ( not .ClientStreaming ) ( not .Paged ) }}
		resp, err := client.{{ .Method }}(ctx, &{{ $inputVar }})
		{{ else if .Paged }}
		iter := client.{{ .Method }}(ctx, &{{ $inputVar }})
		{{ else }}
		resp, err := client.{{ .Method }}(ctx)
		{{ end }}
		{{ if and .ServerStreaming ( not .ClientStreaming ) }}
		var item *{{ .OutputType }}
		for err == nil {
			item, err = resp.Recv()
			fmt.Println(item)
		}

		if err == io.EOF {
			return nil
		}
		{{ else if and .ClientStreaming ( not .ServerStreaming ) }}
		fmt.Println("Client stream open. Close with blank line.")
		scanner := bufio.NewScanner(os.Stdin)
    for err == nil && scanner.Scan() {
				input := scanner.Text()
				if input == "" {
					break
				}
        err = json.Unmarshal([]byte(input), &{{ $inputVar }})
				if err != nil {
					return err
				}
				
				err = resp.Send(&{{ $inputVar }})
    }
    if err := scanner.Err(); err != nil {
        return err
    }
		
		srvResp, err := resp.CloseAndRecv()
		fmt.Println(srvResp)
		{{ else if and .ClientStreaming .ServerStreaming }}
		{{ else if .Paged }}
		var page *{{ .OutputType }}
		for err == nil {
			page, err = iter.Next()

			if err == iterator.Done {
				return nil
			}

			fmt.Println(page)
		}
		{{ else }}
		fmt.Println(resp)
		{{ end }}
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
