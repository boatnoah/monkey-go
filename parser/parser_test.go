package parser

import (
	"fmt"
	"testing"

	"github.com/boatnoah/monkey/ast"
	"github.com/boatnoah/monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 10;
	let y = 10;
	let foobar = 10298371209837;
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() return nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Program does not contain 3 statements, got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Fatalf("s.TokenLiteral is not a 'let' statement. got=%q", s.TokenLiteral())
		return false
	}

	letStatement, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s is not a *ast.LetStatement. got=%T", s)
		return false
	}

	if letStatement.Name.Value != name {
		t.Fatalf("letStatement.Name.Value does not '%s'. got=%s", name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Fatalf("letStatement.Name.TokenLiteral() does not '%s'. got=%s", name, letStatement.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
		return 10;
		return 9;
		return 102398731;
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() return nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Program does not contain 3 statements, got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("Statement is not a return statement. got=%T", returnStmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("return.Stmt.TokenLiteral is not a return statement. got=%s", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	// simplest form of a expression: anotherVar;

	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0].(*ast.ExpressionStatement) is not a expression statement, got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("stmt.Expression is not a identifier, got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s, got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() not %s, got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0].(*ast.ExpressionStatement) is not a expression statement, got=%T", program.Statements[0])
	}

	int, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not a IntegerLiteral, got=%T", stmt.Expression)
	}

	if int.Value != 5 {
		t.Errorf("ident.Value not %d, got=%d", 5, int.Value)
	}

	if int.TokenLiteral() != "5" {
		t.Errorf("int.TokenLiteral() not %s, got=%s", "5", int.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTest := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTest {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program does not have enough statements, got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0].(*ast.ExpressionStatement) is not a expression statement, got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt.Expression is not a IntegerLiteral, got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}

	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.errors
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error %q errors", msg)
	}

	t.FailNow()
}
