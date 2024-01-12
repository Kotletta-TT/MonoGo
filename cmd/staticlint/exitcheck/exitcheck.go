// Package exitcheck implements a static analyzer that checks for an exit in the main function.
package exitcheck

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// MainExitCheckAnalyzer is a static analyzer that checks for an exit in the main function.
var MainExitCheckAnalyzer = &analysis.Analyzer{
	Name: "mainexitcheck",
	Doc:  "Check for exit in main function",
	Run:  run,
}

// run is a Go function that inspects the main function in the main package for an exit call and reports it if found.
//
// It takes a pointer to analysis.Pass as a parameter and returns an interface and an error.
func run(pass *analysis.Pass) (interface{}, error) {
	packageName := ""
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			if decl, ok := n.(*ast.File); ok {
				packageName = decl.Name.Name
				return true
			}

			if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Name.Name == "main" && packageName == "main" {
				for _, stmt := range funcDecl.Body.List {
					if callExpr, ok := stmt.(*ast.ExprStmt); ok {
						if call, ok := callExpr.X.(*ast.CallExpr); ok {
							if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
								if obj, ok := fun.X.(*ast.Ident); ok && obj.Name == "os" && fun.Sel.Name == "Exit" {
									pass.Reportf(callExpr.Pos(), "exit in main function")
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
