package parser

import (
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"unicode"
)

type Expr struct {
	Num   string
	Denom string
}

func (f Expr) String() string {
	return fmt.Sprintf("Expr(Num: %q, Denum: %q)", f.Num, f.Denom)
}

func NewRatFromExpr(f Expr) (*big.Rat, error) {
	if strings.TrimSpace(f.Num) == "" {
		return nil, errors.New("numerator is empty")
	}

	numRatExpr, err := parseExpressionToRat(f.Num)
	if err != nil {
		return nil, err
	}

	if slices.Contains([]string{"", "1"}, strings.TrimSpace(f.Denom)) {
		return numRatExpr, nil
	}

	denRatExpr, err := parseExpressionToRat(f.Denom)
	if err != nil {
		return nil, err
	}

	if denRatExpr.Sign() == 0 {
		return nil, fmt.Errorf("divide by zero")
	}

	return new(big.Rat).Quo(numRatExpr, denRatExpr), nil
}

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	WS
	NUM
	LPAREN
	RPAREN
	ADD
	SUB
	MUL
	QUO
	POWER
	MODULO
)

type Token struct {
	Typ   TokenType
	Value string
}

func (t Token) String() string {
	var typName string
	switch t.Typ {
	case ILLEGAL:
		typName = "ILLEGAL"
	case EOF:
		typName = "EOF"
	case WS:
		typName = "WS"
	case NUM:
		typName = "NUM"
	case LPAREN:
		typName = "LPAREN"
	case RPAREN:
		typName = "RPAREN"
	case ADD:
		typName = "ADD"
	case SUB:
		typName = "SUB"
	case MUL:
		typName = "MUL"
	case QUO:
		typName = "QUO"
	case POWER:
		typName = "POWER"
	case MODULO:
		typName = "MODULO"
	default:
		typName = "UNKNOWN"
	}

	return fmt.Sprintf("Token(%s, %q)", typName, t.Value)
}

type Scanner struct {
	s   []rune
	pos int
}

func NewScanner(s string) *Scanner {
	return &Scanner{s: []rune(s)}
}

func (sc *Scanner) scan() Token {
	if sc.pos >= len(sc.s) {
		return Token{Typ: EOF}
	}
	if unicode.IsSpace(sc.s[sc.pos]) {
		start := sc.pos
		for sc.pos < len(sc.s) && unicode.IsSpace(sc.s[sc.pos]) {
			sc.pos++
		}
		return Token{Typ: WS, Value: string(sc.s[start:sc.pos])}
	}

	if sc.pos+1 < len(sc.s) && sc.s[sc.pos] == '*' && sc.s[sc.pos+1] == '*' {
		sc.pos += 2
		return Token{Typ: POWER, Value: "**"}
	}

	ch := sc.s[sc.pos]
	start := sc.pos

	if unicode.IsDigit(ch) {
		sc.pos++
		for sc.pos < len(sc.s) && unicode.IsDigit(sc.s[sc.pos]) {
			sc.pos++
		}
		if sc.pos < len(sc.s) && sc.s[sc.pos] == '.' {
			if sc.pos+1 >= len(sc.s) || !unicode.IsDigit(sc.s[sc.pos+1]) {
				sc.pos++
				return Token{Typ: ILLEGAL, Value: string(sc.s[start:sc.pos])}
			}
			sc.pos++
			for sc.pos < len(sc.s) && unicode.IsDigit(sc.s[sc.pos]) {
				sc.pos++
			}
		}
		return Token{Typ: NUM, Value: string(sc.s[start:sc.pos])}
	} else if ch == '.' {
		if sc.pos+1 >= len(sc.s) || !unicode.IsDigit(sc.s[sc.pos+1]) {
			sc.pos++
			return Token{Typ: ILLEGAL, Value: string(sc.s[start:sc.pos])}
		}
		sc.pos++
		for sc.pos < len(sc.s) && unicode.IsDigit(sc.s[sc.pos]) {
			sc.pos++
		}
		return Token{Typ: NUM, Value: string(sc.s[start:sc.pos])}
	}

	valStr := string(ch)
	sc.pos++
	switch ch {
	case '(':
		return Token{Typ: LPAREN, Value: valStr}
	case ')':
		return Token{Typ: RPAREN, Value: valStr}
	case '+':
		return Token{Typ: ADD, Value: valStr}
	case '-':
		return Token{Typ: SUB, Value: valStr}
	case '*':
		return Token{Typ: MUL, Value: valStr}
	case '/':
		return Token{Typ: QUO, Value: valStr}
	case '%':
		return Token{Typ: MODULO, Value: valStr}
	}
	return Token{Typ: ILLEGAL, Value: valStr}
}

func scanAll(s string) ([]Token, error) {
	scanner := NewScanner(s)
	var tokens []Token
	for {
		t := scanner.scan()
		if t.Typ == ILLEGAL {
			return nil, fmt.Errorf("illegal token encountered %q near position %d", t.Value, scanner.pos)
		}
		if t.Typ != WS {
			tokens = append(tokens, t)
		}
		if t.Typ == EOF {
			break
		}
	}
	return tokens, nil
}

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) peek() Token {
	if p.pos >= len(p.tokens) {
		return Token{Typ: EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) consume(expected TokenType) (Token, error) {
	t := p.peek()
	if t.Typ == EOF && expected != EOF {
		return t, fmt.Errorf("unexpected EOF")
	}
	if t.Typ != expected {
		return t, fmt.Errorf("unexpected Token")
	}
	p.pos++
	return t, nil
}

func (p *Parser) parseExpression() (*big.Rat, error) {
	lhs, err := p.parseTerm()
	if err != nil {
		return new(big.Rat), err
	}

	for {
		t := p.peek()
		switch t.Typ {
		case ADD:
			p.consume(t.Typ)
			rhs, err := p.parseTerm()
			if err != nil {
				return new(big.Rat), err
			}
			lhs.Add(lhs, rhs)
		case SUB:
			p.consume(t.Typ)
			rhs, err := p.parseTerm()
			if err != nil {
				return new(big.Rat), err
			}
			lhs.Sub(lhs, rhs) // BUG FIX: Add -> Sub
		default:
			return lhs, nil
		}
	}
}

func (p *Parser) parseTerm() (*big.Rat, error) {
	lhs, err := p.parsePower()
	if err != nil {
		return new(big.Rat), err
	}

	for {
		t := p.peek()
		switch t.Typ {
		case MUL:
			p.consume(t.Typ)
			rhs, err := p.parsePower()
			if err != nil {
				return new(big.Rat), err
			}
			lhs.Mul(lhs, rhs)
		case QUO:
			p.consume(t.Typ)
			rhs, err := p.parsePower()
			if err != nil {
				return new(big.Rat), err
			}
			if rhs.Sign() == 0 {
				return new(big.Rat), errors.New("division by zero in expression")
			}
			lhs.Quo(lhs, rhs)
		case MODULO:
			p.consume(t.Typ)
			rhs, err := p.parsePower()
			if err != nil {
				return new(big.Rat), err
			}
			if !lhs.IsInt() || !rhs.IsInt() {
				return new(big.Rat), errors.New("modulo operator requires integer operands") // FIX: Improved error message
			}
			lhsNum := lhs.Num()
			rhsNum := rhs.Num()
			resultNum := new(big.Int).Rem(lhsNum, rhsNum)
			lhs.SetInt(resultNum)
		default:
			return lhs, nil
		}
	}
}

func (p *Parser) parsePower() (*big.Rat, error) {
	lhs, err := p.parseUnary()
	if err != nil {
		return new(big.Rat), err
	}
	if p.peek().Typ == POWER {
		p.consume(POWER)
		rhs, err := p.parsePower()
		if err != nil {
			return new(big.Rat), err
		}
		return calculatePower(lhs, rhs)
	}

	return lhs, nil
}

func (p *Parser) parseUnary() (*big.Rat, error) {
	t := p.peek()
	if t.Typ == ADD {
		p.consume(ADD)
		val, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return val, nil
	}
	if t.Typ == SUB {
		p.consume(SUB)
		val, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return new(big.Rat).Neg(val), nil
	}

	return p.parseAtom()
}

func (p *Parser) parseAtom() (*big.Rat, error) {
	t := p.peek()
	switch t.Typ {
	case NUM:
		p.consume(t.Typ)
		val := new(big.Rat)
		if _, ok := val.SetString(t.Value); !ok {
			return nil, fmt.Errorf("invalid number format for big.Rat: %q", t.Value)
		}
		return val, nil
	case LPAREN:
		p.consume(t.Typ)
		val, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RPAREN); err != nil {
			return nil, fmt.Errorf("missing closing parenthesis: %w", err)
		}
		return val, nil
	default:
		if t.Typ == EOF {
			return nil, errors.New("unexpected EOF. expected a number or parenthesis")
		}
		return nil, fmt.Errorf("unexpected token in atom: %s (expected number or '(')", t)
	}
}

func calculatePower(base, exponent *big.Rat) (*big.Rat, error) {
	if !exponent.IsInt() {
		return nil, fmt.Errorf("exponent must be integer for power operation")
	}
	expVal := exponent.Num()
	if base.Sign() == 0 && expVal.Sign() < 0 {
		return nil, errors.New("math error: 0 raised to a negative power")
	}

	num, den, finalNum, finalDen := base.Num(), base.Denom(), new(big.Int), new(big.Int)
	if expVal.Sign() >= 0 {
		finalNum.Exp(num, expVal, nil)
		finalDen.Exp(den, expVal, nil)
	} else {
		absExpVal := new(big.Int).Abs(expVal)
		finalNum.Exp(den, absExpVal, nil)
		finalDen.Exp(num, absExpVal, nil)
	}
	if finalDen.Sign() == 0 {
		return nil, errors.New("math error: division by zero in power calculation")
	}

	return new(big.Rat).SetFrac(finalNum, finalDen), nil
}

func parseExpressionToRat(s string) (*big.Rat, error) {
	tokens, err := scanAll(s)
	if err != nil {
		return nil, fmt.Errorf("scanAll failed for %q: %w", s, err)
	}

	var hasContent bool
	if len(tokens) > 0 {
		if !(len(tokens) == 1 && tokens[0].Typ == EOF) {
			hasContent = true
		}
	}
	if !hasContent {
		return nil, fmt.Errorf("no evaluatable expression in string: %q", s)
	}

	p := NewParser(tokens)
	ret, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("parseExpression failed for %q: %w", s, err)
	}
	if p.peek().Typ != EOF {
		return nil, fmt.Errorf("unexpected trailling tokens starting with: %s in %q", p.peek(), s)
	}

	return ret, nil
}
