package analyzer

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type Config struct {
	RejectBasic bool // reject basic type return value
}

func New() *analysis.Analyzer {
	return NewFromConfig(&Config{})
}

func NewFromConfig(c *Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:  "deferresult",
		Doc:   "finds potentially unhandled return value of defer statement",
		Run:   c.Run,
		Flags: c.FlagSet(),
	}
}

func (c *Config) FlagSet() flag.FlagSet {
	var fs flag.FlagSet
	fs.BoolVar(&(c.RejectBasic), "nobasic", false, "reject basic values")
	return fs
}

func (c *Config) Run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.DeferStmt:
				retType := pass.TypesInfo.Types[x.Call].Type

				if t, ok := retType.(*types.Tuple); ok && t.Len() == 0 {
					// void
					return true
				}
				if !c.RejectBasic {
					if _, ok := retType.(*types.Basic); ok {
						// basic types
						return true
					}
				}

				diag := analysis.Diagnostic{
					Pos:     x.Defer,
					Message: fmt.Sprintf("unhandled return (%s) in defer statement", retType),
				}
				pass.Report(diag)
			default:
			}

			return true
		})
	}
	return nil, nil
}
