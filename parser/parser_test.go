package parser

import (
	"testing"

	"github.com/weiser/lox/expr"
	"github.com/weiser/lox/scanner"
)

func TestIfStmt(t *testing.T) {
	scanner := scanner.MakeScanner(`if (1>0) { print 1;}`)
	toks := scanner.ScanTokens()
	p := Parser{Tokens: toks}
	stmts, _ := p.Parse()

	if _, ok := stmts[0].(*expr.If); !ok {
		t.Errorf("expected an If statement, got a %v", stmts[0])
	}
}

func TestLogicalExpr(t *testing.T) {
	scanner := scanner.MakeScanner(`var a = null or 1;`)
	toks := scanner.ScanTokens()
	p := Parser{Tokens: toks}
	stmts, _ := p.Parse()

	v, ok := stmts[0].(*expr.Var)
	if !ok {
		t.Errorf("expected an asignment statement, got a %v", stmts[0])
	}
	if _, ok := v.Initializer.(*expr.Logical); !ok {
		t.Errorf("Expected a logical expro, got a, %v", v)
	}
}
