package evaluator

import (
	"github.com/kiki-ki/go-monkey/ast"
	"github.com/kiki-ki/go-monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// 文の評価
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	// 式の評価
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var res object.Object
	for _, stmt := range stmts {
		res = Eval(stmt)
	}
	return res
}
