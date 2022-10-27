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
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	var errFound bool
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				if x.Name.String() == "main" {
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
		if exprStmt, ok = stmt.(*ast.ExprStmt); ok {
			if callExpr, ok = exprStmt.X.(*ast.CallExpr); ok {
				if selectorExpr, ok = callExpr.Fun.(*ast.SelectorExpr); ok {
					if ident, ok = selectorExpr.X.(*ast.Ident); ok {
						if ident.Name == "os" && selectorExpr.Sel.Name == "Exit" {
							return true, selectorExpr.Pos()
						}
					}
				}
			}
		}
	}

	return false, f.Pos()
}
