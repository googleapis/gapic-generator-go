package gencli

import (
	"text/template"
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

	g.addGoFile("main.go")

	g.pt.Reset()
}
