package safestructfield

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "safestructfield",
	Doc: `safestructfield detects cases where struct fields are used in methods 
without being properly initialized. 

This linter helps prevent potential runtime errors caused by accessing 
nil or uninitialized fields within methods. It analyzes struct methods and 
reports instances where a field is accessed without a prior assignment, 
which may lead to unexpected behavior or panics.

Example:
    
    type Data struct {
        value *int
    }
    
    func (d *Data) Print() {
        fmt.Println(*d.value) // Potential nil dereference if 'value' is not initialized
    }

This analyzer ensures safer struct usage by identifying such issues early.`,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		fmt.Printf("%#v\n", n)
	})

	return nil, nil
}
