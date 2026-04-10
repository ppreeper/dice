package dice

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

// Expression represents a parsed expression whose Root is the AST root node.
type Expression struct {
	Root *node
}

// node is an AST node for expressions. Binary operators have Op set and Left/Right
// non-nil. Leaf nodes have Dice or Literal set.
type node struct {
	Op      string
	Left    *node
	Right   *node
	Dice    *ParsedDice
	Literal int
}

// ParseExpression parses an expression supporting + - * / (and 'x' as '*') and
// parentheses. Multiplication/division have higher precedence than addition/subtraction.
// Dice terms are parsed with the existing Parse function; integer literals are supported.
func ParseExpression(s string) (Expression, error) {
	var expr Expression
	s = strings.TrimSpace(s)
	if s == "" {
		return expr, fmt.Errorf("empty expression")
	}
	p := &parser{s: s}
	n, err := p.parseExpression()
	if err != nil {
		return expr, err
	}
	p.skipSpaces()
	if p.pos < len(p.s) {
		return expr, fmt.Errorf("unexpected token after expression: %q", p.s[p.pos:])
	}
	expr.Root = n
	return expr, nil
}

// RollExpression evaluates the expression AST by rolling dice terms and combining
// values using operator precedence. It aggregates per-dice metadata (AllRolls, Rolls,
// RerollsPerformed, TotalRollCalls, Successes) from dice subterms.
func RollExpression(expr Expression, rng *rand.Rand) (RollResult, error) {
	if rng == nil {
		rng = defaultRand
	}
	var combined RollResult
	if expr.Root == nil {
		return combined, fmt.Errorf("empty expression")
	}
	r, val, err := evalNode(expr.Root, rng)
	if err != nil {
		return combined, err
	}
	combined = r
	combined.Total = val
	return combined, nil
}

// evalNode evaluates an AST node and returns aggregated RollResult and computed int value.
func evalNode(n *node, rng *rand.Rand) (RollResult, int, error) {
	var res RollResult
	if n == nil {
		return res, 0, fmt.Errorf("nil node")
	}
	if n.Dice != nil {
		r, err := RollParsed(*n.Dice, rng)
		if err != nil {
			return res, 0, err
		}
		return r, r.Total, nil
	}
	if n.Left == nil && n.Right == nil {
		// literal
		return res, n.Literal, nil
	}
	// binary operator
	leftRes, leftVal, err := evalNode(n.Left, rng)
	if err != nil {
		return res, 0, err
	}
	rightRes, rightVal, err := evalNode(n.Right, rng)
	if err != nil {
		return res, 0, err
	}
	// aggregate metadata
	res.AllRolls = append(res.AllRolls, leftRes.AllRolls...)
	res.AllRolls = append(res.AllRolls, rightRes.AllRolls...)
	res.Rolls = append(res.Rolls, leftRes.Rolls...)
	res.Rolls = append(res.Rolls, rightRes.Rolls...)
	res.RerollsPerformed = leftRes.RerollsPerformed + rightRes.RerollsPerformed
	res.TotalRollCalls = leftRes.TotalRollCalls + rightRes.TotalRollCalls
	res.Successes = leftRes.Successes + rightRes.Successes

	switch n.Op {
	case "+":
		return res, leftVal + rightVal, nil
	case "-":
		return res, leftVal - rightVal, nil
	case "*":
		return res, leftVal * rightVal, nil
	case "/":
		if rightVal == 0 {
			return res, 0, fmt.Errorf("division by zero")
		}
		return res, leftVal / rightVal, nil
	default:
		return res, 0, fmt.Errorf("unsupported operator %q", n.Op)
	}
}

// parser implements a simple recursive-descent parser for expressions.
type parser struct {
	s   string
	pos int
}

func (p *parser) skipSpaces() {
	for p.pos < len(p.s) && p.s[p.pos] == ' ' {
		p.pos++
	}
}

// parseExpression := parseTerm { ('+'|'-') parseTerm }
func (p *parser) parseExpression() (*node, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for {
		p.skipSpaces()
		if p.pos >= len(p.s) {
			break
		}
		ch := p.s[p.pos]
		if ch != '+' && ch != '-' {
			break
		}
		p.pos++
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &node{Op: string(ch), Left: left, Right: right}
	}
	return left, nil
}

// parseTerm := parseFactor { ('*'|'/'|'x') parseFactor }
func (p *parser) parseTerm() (*node, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for {
		p.skipSpaces()
		if p.pos >= len(p.s) {
			break
		}
		ch := p.s[p.pos]
		if ch != '*' && ch != '/' && ch != 'x' {
			break
		}
		p.pos++
		op := string(ch)
		if op == "x" {
			op = "*"
		}
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		left = &node{Op: op, Left: left, Right: right}
	}
	return left, nil
}

// parseFactor := ('+'|'-') parseFactor | parsePrimary
func (p *parser) parseFactor() (*node, error) {
	p.skipSpaces()
	if p.pos < len(p.s) {
		ch := p.s[p.pos]
		if ch == '+' {
			p.pos++
			return p.parseFactor()
		}
		if ch == '-' {
			p.pos++
			// unary minus represented as multiply by -1
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			return &node{Op: "*", Left: &node{Literal: -1}, Right: right}, nil
		}
	}
	return p.parsePrimary()
}

// parsePrimary := '(' parseExpression ')' | dice | literal
func (p *parser) parsePrimary() (*node, error) {
	p.skipSpaces()
	if p.pos >= len(p.s) {
		return nil, fmt.Errorf("unexpected end of expression")
	}
	if p.s[p.pos] == '(' {
		p.pos++
		n, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		p.skipSpaces()
		if p.pos >= len(p.s) || p.s[p.pos] != ')' {
			return nil, fmt.Errorf("missing closing parenthesis")
		}
		p.pos++
		return n, nil
	}
	// read token until next top-level operator or end or closing parenthesis
	start := p.pos
	depth := 0
	for p.pos < len(p.s) {
		ch := p.s[p.pos]
		if ch == '(' {
			depth++
		} else if ch == ')' {
			if depth == 0 {
				break
			}
			depth--
		} else if depth == 0 && (ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == 'x') {
			break
		}
		p.pos++
	}
	token := strings.TrimSpace(p.s[start:p.pos])
	if token == "" {
		return nil, fmt.Errorf("empty term in expression")
	}
	// dice term if contains 'd' or ends with 'F'
	if strings.Contains(token, "d") || strings.HasSuffix(token, "F") {
		pd, err := Parse(token)
		if err != nil {
			return nil, fmt.Errorf("parse dice term %q: %w", token, err)
		}
		return &node{Dice: &pd}, nil
	}
	// integer literal
	v, err := strconv.Atoi(token)
	if err != nil {
		return nil, fmt.Errorf("invalid literal %q: %w", token, err)
	}
	return &node{Literal: v}, nil
}
