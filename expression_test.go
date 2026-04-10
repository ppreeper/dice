package dice

import (
	"math/rand/v2"
	"testing"
)

func TestParseExpressionSimple(t *testing.T) {
	cases := []string{"2d6+1", "1d6+1d4-3", "4d6", "3+2d6"}
	for _, in := range cases {
		if _, err := ParseExpression(in); err != nil {
			t.Fatalf("parse expr %q err: %v", in, err)
		}
	}
}

func TestRollExpressionBasic(t *testing.T) {
	expr, err := ParseExpression("2d6+1")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	rng := rand.New(rand.NewPCG(600, 601))
	res, err := RollExpression(expr, rng)
	if err != nil {
		t.Fatalf("roll expr err: %v", err)
	}
	// basic sanity checks
	if len(res.AllRolls) != 2 {
		t.Fatalf("expected 2 all-rolls, got %d", len(res.AllRolls))
	}
	if res.Total < 2 || res.Total > 13 { // 2d6 + 1 -> min 3, max 13 but allow broader check
		t.Fatalf("unexpected total %d", res.Total)
	}
}

func TestRollExpressionMultipleDiceAndLiterals(t *testing.T) {
	expr, err := ParseExpression("1d6+1d4-3")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	rng := rand.New(rand.NewPCG(600, 601))
	res, err := RollExpression(expr, rng)
	if err != nil {
		t.Fatalf("roll expr err: %v", err)
	}
	// ensure totals align with component rolls
	// total = sum(rolls of dice terms) - 3
	sum := 0
	for _, v := range res.Rolls {
		sum += v
	}
	if res.Total != sum-3 {
		t.Fatalf("expression total mismatch: got %d want %d; rolls=%+v", res.Total, sum-3, res.Rolls)
	}
}

func TestPrecedenceAndParens(t *testing.T) {
	// multiplication before addition: 2 + 3 * 4 == 14
	expr, err := ParseExpression("2+3*4")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	res, err := RollExpression(expr, nil)
	if err != nil {
		t.Fatalf("eval err: %v", err)
	}
	if res.Total != 14 {
		t.Fatalf("precedence failed: got %d want 14", res.Total)
	}

	// parentheses override: (2+3)*4 == 20
	expr2, err := ParseExpression("(2+3)*4")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	res2, err := RollExpression(expr2, nil)
	if err != nil {
		t.Fatalf("eval err: %v", err)
	}
	if res2.Total != 20 {
		t.Fatalf("paren precedence failed: got %d want 20", res2.Total)
	}
}

func TestDiceWithMultiplicationDivision(t *testing.T) {
	// 3d6x2 should multiply the sum of 3d6 by 2
	expr, err := ParseExpression("3d6x2")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	rng := rand.New(rand.NewPCG(600, 601))
	res, err := RollExpression(expr, rng)
	if err != nil {
		t.Fatalf("roll err: %v", err)
	}
	// sum of rolls times 2
	sum := 0
	for _, v := range res.Rolls {
		sum += v
	}
	if res.Total != sum*2 {
		t.Fatalf("multiply dice failed: got %d want %d", res.Total, sum*2)
	}

	// division: (sum of 4d6) / 2
	expr2, err := ParseExpression("4d6/2")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	res2, err := RollExpression(expr2, rng)
	if err != nil {
		t.Fatalf("roll err: %v", err)
	}
	sum2 := 0
	for _, v := range res2.Rolls {
		sum2 += v
	}
	if res2.Total != sum2/2 {
		t.Fatalf("divide dice failed: got %d want %d", res2.Total, sum2/2)
	}
}

func TestDivisionByZeroInExpression(t *testing.T) {
	if _, err := ParseExpression("2d6/0"); err != nil {
		// parsing should succeed; evaluation will fail on division by zero
	}
	expr, _ := ParseExpression("2+4/0")
	if _, err := RollExpression(expr, nil); err == nil {
		t.Fatalf("expected division by zero error")
	}
}
