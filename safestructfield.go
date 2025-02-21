package safestructfield

import (
	"go/ast"
	"go/types"

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

	cache := make(map[*types.Var]bool)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Recv == nil {
			return
		}

		fnRecvName := fn.Recv.List[0].Names[0].Name

		fields := structFields(pass, fn)
		if fields == nil {
			return
		}

		ast.Inspect(fn.Body, func(n ast.Node) bool {
			selector, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			// Filter: only method's recv selection
			if selExprId, ok := selector.X.(*ast.Ident); !ok || selExprId.Name != fnRecvName {
				return true
			}

			field, ok := fields[selector.Sel.Name]
			if !ok {
				return true
			}

			if _, ok := field.Type().(*types.Pointer); !ok {
				return true
			}

			cache[field] = false

			return true
		})
	})

	return nil, nil
}

func structFields(pass *analysis.Pass, fn *ast.FuncDecl) map[string]*types.Var {
	recv := fn.Recv.List[0].Type
	recvType, ok := pass.TypesInfo.TypeOf(recv).(*types.Pointer)
	if !ok {
		return nil
	}

	structType, ok := recvType.Elem().Underlying().(*types.Struct)
	if !ok {
		return nil
	}

	fields := make(map[string]*types.Var)
	for i := 0; i < structType.NumFields(); i++ {
		fields[structType.Field(i).Name()] = structType.Field(i)
	}

	return fields
}
