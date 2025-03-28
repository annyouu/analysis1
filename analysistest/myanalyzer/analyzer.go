package myanalyzer

import (
    "go/ast"
    "go/types"
    "golang.org/x/tools/go/analysis"
    "golang.org/x/tools/go/analysis/passes/inspect"
    "golang.org/x/tools/go/ast/inspector"
)

// Analyzerは型アサーションの問題を検出する
var Analyzer = &analysis.Analyzer{
	Name: "fourcetypeassert",
	Doc: "Detects unsafe type assertions",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	// inspect.Analyzerを使用する
	ins, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	// 型アサーションのノードを取得
	nodeFilter := []ast.Node{
		(*ast.TypeAssertExpr)(nil),
	}

	ins.Preorder(nodeFilter, func(n ast.Node) {
		typeAssert, _ := n.(*ast.TypeAssertExpr)

		// nilチェック
		if typeAssert.X == nil {
			return
		}

		// `switch x := y.(type)` の場合は警告を出さない
		if isInsideTypeSwitch(typeAssert, pass) {
			return
		}

		// 型情報が取得できるかチェック
		typ, ok := pass.TypesInfo.Types[typeAssert.X]
		if !ok || typ.Type == nil {
			return
		}

		// typeAssert.Xがinterface{}型かどうかを判定
		if _, ok := typ.Type.(*types.Interface); ok {
			pass.Reportf(typeAssert.Pos(), "unsafe type assertion")
		}
	})
	return nil, nil
}

// 型アサーションが `switch x := y.(type)` の中で使われているかを判定
func isInsideTypeSwitch(typeAssert *ast.TypeAssertExpr, pass *analysis.Pass) bool {
	found := false
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if typeSwitch, ok := n.(*ast.TypeSwitchStmt); ok {
				// `Assign` の中身を解析
				if exprStmt, ok := typeSwitch.Assign.(*ast.ExprStmt); ok {
					if expr, ok := exprStmt.X.(*ast.TypeAssertExpr); ok && expr == typeAssert {
						found = true
						return false // 早期終了
					}
				}
				if assignStmt, ok := typeSwitch.Assign.(*ast.AssignStmt); ok {
					// `AssignStmt` の右辺の式をチェック
					for _, rhs := range assignStmt.Rhs {
						if expr, ok := rhs.(*ast.TypeAssertExpr); ok && expr == typeAssert {
							found = true
							return false // 早期終了
						}
					}
				}
			}
			return true
		})
		if found {
			return true
		}
	}
	return false
}
