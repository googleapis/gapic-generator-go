package gencli

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

const (
	// ExpectedParams is the number of expected plugin parameters
	ExpectedParams = 2

	// ShortDescMax is the maximum length accepted for
	// the Short usage docs
	ShortDescMax = 50
)

func parseParameters(params *string) (rootDir string, gapicDir string, err error) {
	if params == nil {
		err = fmt.Errorf("Missing required parameters. See usage")
		return
	}

	split := strings.Split(*params, ",")
	if len(split) != ExpectedParams {
		err = fmt.Errorf("Improper number of parameters. Got %d, require %d. See usage", len(split), ExpectedParams)
		return
	}

	for _, str := range split {
		sepNdx := strings.Index(str, ":")
		if sepNdx == -1 {
			err = fmt.Errorf("Unknown parameter: %s", str)
			return
		}

		switch str[:sepNdx] {
		case "gapic":
			gapicDir = str[sepNdx+1:]
		case "root":
			rootDir = str[sepNdx+1:]
		default:
			err = fmt.Errorf("Unknown parameter: %s", str)
		}
	}

	return
}

func toShortUsage(cmt string) string {
	if len(cmt) > ShortDescMax {
		sep := strings.LastIndex(cmt[:ShortDescMax], " ")
		if sep == -1 {
			sep = ShortDescMax
		}
		cmt = cmt[:sep] + "..."
	}

	return cmt
}

func sanitizeComment(cmt string) string {
	cmt = strings.Replace(cmt, "\\", `\\`, -1)
	cmt = strings.Replace(cmt, "\n", " ", -1)
	cmt = strings.TrimSpace(cmt)
	return cmt
}

func strContains(a []string, s string) bool {
	for _, as := range a {
		if as == s {
			return true
		}
	}
	return false
}

func copyImports(from, to map[string]*pbinfo.ImportSpec) {
	for _, val := range from {
		putImport(to, val)
	}
}

func putImport(imports map[string]*pbinfo.ImportSpec, pkg *pbinfo.ImportSpec) {
	if _, ok := imports[pkg.Path]; !ok {
		imports[pkg.Path] = pkg
	}
}

func parseMessageName(field *descriptor.FieldDescriptorProto, msg *descriptor.DescriptorProto) (name string) {
	t := field.GetTypeName()
	last := strings.LastIndex(t, ".")
	name = t[last+1:]

	// check if it is a nested type
	if strings.Contains(t, msg.GetName()) {
		pre := t[:last]
		parent := pre[strings.LastIndex(pre, ".")+1:]
		name = parent + "_" + name
	}

	return
}

func title(name string) string {
	split := strings.Split(name, "_")
	for ndx, tkn := range split {
		split[ndx] = strings.Title(tkn)
	}

	return strings.Join(split, "")
}

func isNestedOneOf(in, cur string) bool {
	in = in[strings.LastIndex(in, ".")+1:]

	return in != cur
}
