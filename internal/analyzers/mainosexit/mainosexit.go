package mainosexit

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "mainosexit",
	Doc:  Doc,
	Run:  run,
}

const Doc = `detect whether the os.Exit function is being used in the main function

The mainosexit analysis reports whether a main function explicitly uses os.Exit function.`

func run(pass *analysis.Pass) (interface{}, error) {
	const (
		packName string = "main"
		funcName string = "main"
	)
	if pass.Pkg.Name() != packName {
		return nil, nil
	}

	var errFound bool
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.Package:
			case *ast.FuncDecl:
				if x.Name.String() == funcName {
					if spot, pos := doesFuncCallExit(x); spot {
						pass.Reportf(pos, "the main function explicitly calls for os.Exit()")
						errFound = true
						return false
					}
				}
			}
			return true
		})

		if errFound {
			break
		}
	}
	return nil, nil
}

func doesFuncCallExit(f *ast.FuncDecl) (bool, token.Pos) {
	var ok bool
	var exprStmt *ast.ExprStmt
	var callExpr *ast.CallExpr
	var selectorExpr *ast.SelectorExpr
	var ident *ast.Ident

	for _, stmt := range f.Body.List {
		if exprStmt, ok = stmt.(*ast.ExprStmt); !ok {
			continue
		}

		if callExpr, ok = exprStmt.X.(*ast.CallExpr); !ok {
			continue
		}

		if selectorExpr, ok = callExpr.Fun.(*ast.SelectorExpr); !ok {
			continue
		}

		if ident, ok = selectorExpr.X.(*ast.Ident); !ok {
			continue
		}

		if ident.Name == "os" && selectorExpr.Sel.Name == "Exit" {
			return true, selectorExpr.Pos()
		}
	}

	return false, f.Pos()
}
