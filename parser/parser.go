package parser

import (
	"fmt"

	"github.com/weiser/lox/expr"
	"github.com/weiser/lox/token"
)

type ParserError struct {
	Token token.Token
	Msg   string
}

func (pe *ParserError) Error() string {
	return fmt.Sprintf("error on %v: %v", pe.Token, pe.Msg)
}

type Parser struct {
	Tokens     []token.Token
	ParsingErr *ParserError
	Current    int
}

func (p *Parser) Parse() ([]expr.StmtInterface, error) {
	stmts := make([]expr.StmtInterface, 0)
	for !p.isAtEnd() {
		stmts = append(stmts, p.Declaration())
	}
	return stmts, nil
}

func (p *Parser) Declaration() expr.StmtInterface {
	defer func() {
		if err := recover(); err != nil {
			v, ok := err.(*ParserError)
			if ok {
				p.ParsingErr = v
				p.synchronize()
			} else {
				// any non-parsererror we will barf on
				panic(v)
			}
		}
	}()
	if p.match(token.CLASS) {
		return p.ClassDeclaration()
	}
	if p.match(token.FUN) {
		return p.Function("function")
	}
	if p.match(token.VAR) {
		return p.VarDeclaration()
	}
	return p.Statement()
}

func (p *Parser) ClassDeclaration() expr.StmtInterface {
	name, err := p.consume(token.IDENTIFIER, "expect class name")
	if err != nil {
		panic(err)
	}
	_, errlb := p.consume(token.LEFT_BRACE, "expect '{' before class body")
	if errlb != nil {
		panic(errlb)
	}

	methods := make([]expr.StmtInterface, 0)
	for !p.checkType(token.RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.Function("method"))
	}

	_, errrb := p.consume(token.RIGHT_BRACE, "expect '}' after class body")
	if errrb != nil {
		panic(errrb)
	}

	return &expr.Class{Name: name, Methods: methods}

}

func (p *Parser) Function(kind string) expr.StmtInterface {
	name, err := p.consume(token.IDENTIFIER, fmt.Sprintf("Expect %v name", kind))
	if err != nil {
		panic(err)
	}
	_, perr := p.consume(token.LEFT_PAREN, fmt.Sprintf("Expect '(' after %v name", kind))
	if perr != nil {
		panic(perr)
	}

	parameters := make([]token.Token, 0)
	if !p.checkType(token.RIGHT_PAREN) {
		// consume 1 or more params
		param, iderr := p.consume(token.IDENTIFIER, "Expect parameter name")
		if iderr != nil {
			panic(iderr)
		}
		parameters = append(parameters, param)
		for p.match(token.COMMA) {
			param, iderr = p.consume(token.IDENTIFIER, "Expect parameter name")
			if iderr != nil {
				panic(iderr)
			}
			parameters = append(parameters, param)
			if len(parameters) >= 255 {
				panic(ParserError{Msg: "cannot have more than 255 args in function call"})
			}
		}
	}
	_, rperr := p.consume(token.RIGHT_PAREN, "Expect ')' after parameters")
	if rperr != nil {
		panic(rperr)
	}

	_, lberr := p.consume(token.LEFT_BRACE, fmt.Sprintf("Expect '{' before %v body.", kind))
	if lberr != nil {
		panic(lberr)
	}

	body := p.BlockStatement()
	fmt.Println("parsed function")
	return &expr.Function{Name: name, Params: parameters, Body: body}
}

func (p *Parser) VarDeclaration() expr.StmtInterface {
	name, _ := p.consume(token.IDENTIFIER, "Expected variable name")

	var initializer expr.ExprInterface
	if p.match(token.EQUAL) {
		initializer = p.Expression()
	}

	p.consume(token.SEMICOLON, "expected ';' after variable declaration")
	return &expr.Var{Name: name, Initializer: initializer}

}

func (p *Parser) Statement() expr.StmtInterface {
	if p.match(token.BREAK) {
		return p.BreakStatement()
	}
	if p.match(token.RETURN) {
		return p.ReturnStatement()
	}
	if p.match(token.FOR) {
		return p.ForStatement()
	}
	if p.match(token.IF) {
		return p.IfStatement()
	}
	if p.match(token.PRINT) {
		return p.PrintStatement()
	}
	if p.match(token.WHILE) {
		return p.WhileStatement()
	}
	if p.match(token.LEFT_BRACE) {
		return &expr.Block{Statements: p.BlockStatement()}
	}
	return p.ExpressionStatement()
}

func (p *Parser) ReturnStatement() expr.StmtInterface {
	keywrd := p.previous()
	var value expr.ExprInterface
	if !p.checkType(token.SEMICOLON) {
		value = p.Expression()
	}

	p.consume(token.SEMICOLON, "Expect ';' after return value")
	return &expr.Return{Keyword: keywrd, Value: value}
}

func (p *Parser) BreakStatement() expr.StmtInterface {
	breakStmt := expr.Expression{Expression: &expr.Literal{Value: token.BREAK}}
	_, err := p.consume(token.SEMICOLON, "Expect ';' after 'break'")
	if err != nil {
		panic(err)
	}
	return &breakStmt
}

func (p *Parser) ForStatement() expr.StmtInterface {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'for'")
	if err != nil {
		panic(err)
	}

	var initializer expr.StmtInterface
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer = p.VarDeclaration()
	} else {
		initializer = p.ExpressionStatement()
	}

	var condition expr.ExprInterface
	if !p.checkType(token.SEMICOLON) {
		condition = p.Expression()
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after for loop condition")
	if err != nil {
		panic(err)
	}

	var increment expr.ExprInterface
	if !p.checkType(token.RIGHT_PAREN) {
		increment = p.Expression()
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after for loop increment clause")
	if err != nil {
		panic(err)
	}

	body := p.Statement()

	if increment != nil {
		body = &expr.Block{Statements: []expr.StmtInterface{body, &expr.Expression{Expression: increment}}}
	}

	if condition == nil {
		condition = &expr.Literal{Value: true}
	}
	body = &expr.While{Condition: condition, Body: body}

	if initializer != nil {
		body = &expr.Block{Statements: []expr.StmtInterface{initializer, body}}
	}

	return body

}

func (p *Parser) WhileStatement() expr.StmtInterface {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'while'")
	if err != nil {
		panic(err)
	}
	condition := p.Expression()
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after condition in 'while'")
	if err != nil {
		panic(err)
	}
	body := p.Statement()

	return &expr.While{Condition: condition, Body: body}
}

func (p *Parser) IfStatement() expr.StmtInterface {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'if'")
	condition := p.Expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after if condition")
	thenBranch := p.Statement()
	var elseBranch expr.StmtInterface
	if p.match(token.ELSE) {
		elseBranch = p.Statement()
	}

	return &expr.If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}

}

func (p *Parser) BlockStatement() []expr.StmtInterface {
	stmts := make([]expr.StmtInterface, 0)
	for !p.checkType(token.RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.Declaration())
	}
	p.consume(token.RIGHT_BRACE, "expected '}' after block")
	return stmts
}

func (p *Parser) PrintStatement() expr.StmtInterface {
	value := p.Expression()
	_, err := p.consume(token.SEMICOLON, "expect ; after value")
	if err != nil {
		panic(err)
	}
	return &expr.Print{Expression: value}
}

func (p *Parser) ExpressionStatement() expr.StmtInterface {
	value := p.Expression()
	_, err := p.consume(token.SEMICOLON, "expect ; after value")
	if err != nil {
		panic(err)
	}
	return &expr.Expression{Expression: value}
}

func (p *Parser) Expression() expr.ExprInterface {
	return p.Assignment()
}

func (p *Parser) Assignment() expr.ExprInterface {
	exp := p.Or()
	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.Assignment()

		if e, ok := exp.(*expr.Variable); ok {
			identifier := e.Name
			return &expr.Assign{Name: identifier, Value: value}
		} else if get, ok1 := exp.(*expr.Get); ok1 {
			return &expr.Set{Object: get.Object, Name: get.Name, Value: value}
		}

		panic(MakeParserError(equals, "Invalid assignment target"))
	}

	return exp
}

func (p *Parser) Or() expr.ExprInterface {
	exp := p.And()

	for p.match(token.OR) {
		operator := p.previous()
		right := p.And()
		exp = &expr.Logical{Left: exp, Operator: operator, Right: right}
	}
	return exp
}

func (p *Parser) And() expr.ExprInterface {
	exp := p.Equality()

	for p.match(token.AND) {
		operator := p.previous()
		right := p.And()
		exp = &expr.Logical{Left: exp, Operator: operator, Right: right}
	}

	return exp
}

func (p *Parser) Equality() expr.ExprInterface {
	exp := p.Comparison()
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.Comparison()
		exp = &expr.Binary{Left: exp, Operator: operator, Right: right}
	}
	return exp
}

func (p *Parser) Comparison() expr.ExprInterface {
	exp := p.Term()
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.Term()
		exp = &expr.Binary{Right: right, Operator: operator, Left: exp}
	}
	return exp
}

func (p *Parser) Term() expr.ExprInterface {
	exp := p.Factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.Factor()
		exp = &expr.Binary{Left: exp, Operator: operator, Right: right}
	}
	return exp
}

func (p *Parser) Factor() expr.ExprInterface {
	exp := p.Unary()
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.Unary()
		exp = &expr.Binary{Left: exp, Operator: operator, Right: right}
	}
	return exp
}

func (p *Parser) Unary() expr.ExprInterface {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.Unary()
		return &expr.Unary{Operator: operator, Right: right}
	}

	return p.Call()
}

func (p *Parser) Call() expr.ExprInterface {
	exp := p.Primary()

	for {
		if p.match(token.LEFT_PAREN) {
			exp = p.FinishCall(exp)
		} else if p.match(token.DOT) {
			name, err := p.consume(token.IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				panic(err)
			}
			exp = &expr.Get{Object: exp, Name: name}
		} else {
			break
		}
	}

	return exp
}

func (p *Parser) FinishCall(callee expr.ExprInterface) expr.ExprInterface {
	args := make([]expr.ExprInterface, 0)
	if !p.checkType(token.RIGHT_PAREN) {
		args = append(args, p.Expression())
		for p.match(token.COMMA) {
			if len(args) >= 255 {
				panic(ParserError{Msg: "cannot have more than 255 args in function call"})
			}
			args = append(args, p.Expression())
		}
	}

	paren, _ := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments")

	return &expr.Call{Callee: callee, Paren: paren, Arguments: args}
}

func (p *Parser) Primary() expr.ExprInterface {
	if p.match(token.FALSE) {
		return &expr.Literal{Value: false}
	}
	if p.match(token.TRUE) {
		return &expr.Literal{Value: true}
	}
	if p.match(token.NIL) {
		return &expr.Literal{Value: nil}
	}
	if p.match(token.NUMBER, token.STRING) {
		return &expr.Literal{Value: p.previous().Literal}
	}

	var e expr.ExprInterface
	if p.match(token.LEFT_PAREN) {
		exp := p.Expression()
		_, err := p.consume(token.RIGHT_PAREN, "Expect ')' after expression")
		if err == nil {
			e = &expr.Grouping{Expression: exp}
			return e
		} else {
			panic(err)
		}
	}
	if p.match(token.IDENTIFIER) {
		return &expr.Variable{Name: p.previous()}
	}
	err := MakeParserError(p.peek(), "expected expression")
	panic(err)

}

func (p *Parser) consume(tokenType token.TType, err string) (token.Token, error) {
	if p.checkType(tokenType) {
		return p.advance(), nil
	}

	return token.Token{}, MakeParserError(p.peek(), err)
}

func MakeParserError(tok token.Token, err string) error {
	return &ParserError{Token: tok, Msg: err}
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().TokenType == token.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case token.CLASS, token.FOR, token.FUN, token.IF, token.PRINT, token.RETURN, token.VAR, token.WHILE:
			return
		}
		p.advance()
	}
}

func (p *Parser) match(tokenTypes ...token.TType) bool {
	for _, typ := range tokenTypes {
		if p.checkType(typ) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) checkType(typ token.TType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == typ
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.Current += 1
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() token.Token {
	return p.Tokens[p.Current-1]
}
