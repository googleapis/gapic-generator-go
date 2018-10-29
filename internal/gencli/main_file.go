package gencli

import (
	"text/template"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

const (
	// MainTemplate is the template string for a generated main.go
	MainTemplate = `package main
func main() {
	Execute()
}
`
)

func (g *gcli) genMainFile() {
	g.pt.Reset()
	template.Must(template.New("main").Parse(MainTemplate)).Execute(g.pt.Writer(), nil)
	g.response.File = append(g.response.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String("main.go"),
		Content: proto.String(g.pt.String()),
	})

	g.pt.Reset()
}
