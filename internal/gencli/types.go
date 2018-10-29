package gencli

// Command intermediate representation of an RPC as a CLI command
type Command struct {
	Service      string
	Method       string
	MethodCmd    string
	InputMessage string
	InputFlags   string
	Flags        []*Flag
	ShortDesc    string
	LongDesc     string
}
