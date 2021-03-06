package scanner

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/weiser/lox/token"
)

type Scanner struct {
	Source               string
	Tokens               []token.Token
	Errors               []Error
	Start, Current, Line int
}

type Error struct {
	Source               string
	Message              string
	Start, Current, Line int
}

func (e Error) String() string {
	return fmt.Sprintf("source=%v, Start=%v, current=%v, line=%v", e.Source, e.Start, e.Current, e.Line)
}

func MakeScanner(src string) Scanner {
	return Scanner{Source: src, Tokens: make([]token.Token, 0), Start: 0, Current: 0, Line: 1}
}

func (s *Scanner) ScanTokens() []token.Token {
	for s.Current < len(s.Source) {
		s.Start = s.Current
		if !s.isAtEnd() {
			s.scanToken()
		}
	}

	s.Tokens = append(s.Tokens, token.Token{TokenType: token.EOF, Lexeme: "", Literal: nil, Line: s.Line})
	return s.Tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)
	case '!':
		var ntt token.TType
		if s.match('=') {
			ntt = token.BANG_EQUAL
		} else {
			ntt = token.BANG
		}
		s.addToken(ntt)
	case '=':
		var ntt token.TType
		if s.match('=') {
			ntt = token.EQUAL_EQUAL
		} else {
			ntt = token.EQUAL
		}
		s.addToken(ntt)
	case '<':
		var ntt token.TType
		if s.match('=') {
			ntt = token.LESS_EQUAL
		} else {
			ntt = token.LESS
		}
		s.addToken(ntt)

	case '>':
		var ntt token.TType
		if s.match('=') {
			ntt = token.GREATER_EQUAL
		} else {
			ntt = token.GREATER
		}
		s.addToken(ntt)
	case '/':
		if s.match('/') {
			// single line comments
			for ; s.peek() != '\n' && !s.isAtEnd(); s.advance() {
			}
		} else if s.match('*') {
			// multi-line comments, e.g.
			// ```
			/*
				line1
				line2
			*/
			// ```
			for ; !s.isAtEnd() && (s.peek() != '*' && s.peekNext() != '/'); s.advance() {
			}
			//skip past lass '/'
			s.advance()
			s.advance()
		} else {
			s.addToken(token.SLASH)
		}
	case ' ', '\r', '\t':
		// ignore non-\n whitespace
	case '\n':
		s.Line += 1
	case '"':
		s.string()
	default:
		if unicode.IsDigit(rune(s.Source[s.Current-1])) {
			s.number()
		} else if unicode.IsLetter(rune(s.Source[s.Current-1])) {
			s.identifier()
		} else {
			s.Errors = append(s.Errors, Error{Source: s.Source[s.Start:s.Current], Line: s.Line, Start: s.Start, Current: s.Current, Message: fmt.Sprintf("unknown token: %v", s.Source[s.Start:s.Current])})
			fmt.Println("Error at line: ", s.Line, s.Source[s.Start:s.Current])
		}
	}
}

var keywords = map[string]token.TType{
	"and":    token.AND,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"for":    token.FOR,
	"fun":    token.FUN,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
	"break":  token.BREAK,
}

func (s *Scanner) identifier() {
	for ; s.isAlphanumeric(s.peek()); s.advance() {
	}
	text := s.Source[s.Start:s.Current]
	toktype, exists := keywords[text]
	if !exists {
		toktype = token.IDENTIFIER
	}
	s.addTokenWithObj(toktype, text)
}

func (s *Scanner) isAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
func (s *Scanner) number() {
	for ; unicode.IsDigit(s.peek()); s.advance() {
	}
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.advance()
		for ; unicode.IsDigit(s.peek()); s.advance() {
		}
	}
	if val, err := strconv.ParseFloat(s.Source[s.Start:s.Current], 64); err == nil {
		s.addTokenWithObj(token.NUMBER, val)
	} else {
		s.Errors = append(s.Errors, Error{Source: s.Source[s.Start:s.Current], Line: s.Line, Start: s.Start, Current: s.Current, Message: fmt.Sprintf("bad number: %v", s.Source[s.Start:s.Current])})
	}
}

func (s *Scanner) peekNext() rune {
	if s.Current+1 >= len(s.Source) {
		return 0
	}
	return rune(s.Source[s.Current])
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.Line += 1
		}
		s.advance()
	}
	if s.isAtEnd() {
		s.Errors = append(s.Errors, Error{Source: s.Source[s.Start:s.Current], Line: s.Line, Start: s.Start, Current: s.Current, Message: "unterminated string"})
		return

	}
	// the closing '"'
	s.advance()
	val := s.Source[s.Start+1 : s.Current-1]
	s.addTokenWithObj(token.STRING, val)
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return rune(s.Source[s.Current])
}

func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if string(s.Source[s.Current]) != string(expected) {
		return false
	}
	s.Current += 1
	return true
}

// TODO: returning a byte here is OK for now, but if we need to allow multi-type utf chars, it should be a rune.
func (s *Scanner) advance() byte {
	b := s.Source[s.Current]
	s.Current += 1
	return b
}

func (s *Scanner) addToken(tok token.TType) {
	s.addTokenWithObj(tok, nil)
}

func (s *Scanner) addTokenWithObj(tok token.TType, obj interface{}) {
	text := s.Source[s.Start:s.Current]
	s.Tokens = append(s.Tokens, token.MakeToken(tok, text, obj, s.Line))
}
