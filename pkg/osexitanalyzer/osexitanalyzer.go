// Package analyzer for scan os.Exit expression in main function
package osexitanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer - instance for multichecker
var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for usage os.Exit in main function",
	Run:  run,
}

func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

func run(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					return true
				}
				return false
			case *ast.ExprStmt:
				if call, ok := x.X.(*ast.CallExpr); ok {
					if isPkgDot(call.Fun, "os", "Exit") {
						pass.Reportf(x.Pos(), "preventing call os.Exit in main function")
					}
				}
				return true
			default:
				return true
			}

		})
	}
	return nil, nil
}
