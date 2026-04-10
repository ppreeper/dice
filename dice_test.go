package dice

import (
	"math/rand/v2"
	"testing"
)

// small, deterministic RNG for tests
var testRNG = rand.New(rand.NewPCG(600, 601))

// fixedRand is provided for deterministic behavior in examples/tests if needed
fixedRand = rand.New(rand.NewPCG(600, 601))

var parseTests = []struct {
	in      string
	want    ParsedDice
	wantErr bool
}{
	{"d6", ParsedDice{Count: 1, Sides: 6, Type: "d", ModFunc: "", ModVal: 0}, false},
	{"1d6", ParsedDice{Count: 1, Sides: 6, Type: "d"}, false},
	{"2d8+3", ParsedDice{Count: 2, Sides: 8, Type: "d", ModFunc: "+", ModVal: 3}, false},
	{"4F", ParsedDice{Count: 4, Sides: 3, Type: "F"}, false},
	{"d%", ParsedDice{Count: 1, Sides: 100, Type: "d"}, false},
	{"", ParsedDice{}, true},
	{"xd6", ParsedDice{}, true},
	{"2d", ParsedDice{}, true},
	{"2d6/0", ParsedDice{Count: 2, Sides: 6, Type: "d", ModFunc: "/", ModVal: 0}, false},
	{"4d6kh3", ParsedDice{Count: 4, Sides: 6, Type: "d", KeepDropAction: "k", KeepDropWhich: "h", KeepDropCount: 3}, false},
	{"4d6!", ParsedDice{Count: 4, Sides: 6, Type: "d", Explode: true}, false},
	{"4d6r1", ParsedDice{Count: 4, Sides: 6, Type: "d", RerollVal: 1, RerollOnce: false}, false},
	{"3d6ro1", ParsedDice{Count: 3, Sides: 6, Type: "d", RerollVal: 1, RerollOnce: true}, false},
	{"10d10>=8", ParsedDice{Count: 10, Sides: 10, Type: "d", SuccessOp: ">=", SuccessVal: 8}, false},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		pd, err := Parse(tt.in)
		if (err != nil) != tt.wantErr {
			t.Fatalf("Parse(%q) err = %v, wantErr %v", tt.in, err, tt.wantErr)
		}
		if tt.wantErr {
			continue
		}
		if pd.Count != tt.want.Count || pd.Sides != tt.want.Sides || pd.Type != tt.want.Type || pd.ModFunc != tt.want.ModFunc || pd.ModVal != tt.want.ModVal {
			t.Fatalf("Parse(%q) = %+v, want %+v", tt.in, pd, tt.want)
		}
	}
}

func TestRollParsedSimple(t *testing.T) {
	pd := ParsedDice{Count: 1, Sides: 6, Type: "d"}
	res, err := RollParsed(pd, testRNG)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rolls) != 1 {
		t.Fatalf("expected 1 roll, got %d", len(res.Rolls))
	}
	if res.Total != res.Rolls[0] {
		t.Fatalf("expected total %d == roll %d", res.Total, res.Rolls[0])
	}
}

func TestFateDice(t *testing.T) {
	pd := ParsedDice{Count: 4, Type: "F", Sides: 3}
	res, err := RollParsed(pd, testRNG)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rolls) != 4 {
		t.Fatalf("expected 4 rolls, got %d", len(res.Rolls))
	}
}

func TestModifiers(t *testing.T) {
	pd := ParsedDice{Count: 2, Sides: 6, Type: "d", ModFunc: "+", ModVal: 3}
	res, err := RollParsed(pd, testRNG)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// total should be sum(rolls) + 3
	sum := 0
	for _, v := range res.Rolls {
		sum += v
	}
	if res.Total != sum+3 {
		t.Fatalf("expected total %d, got %d", sum+3, res.Total)
	}
}

func TestDivisionByZero(t *testing.T) {
	pd := ParsedDice{Count: 2, Sides: 6, Type: "d", ModFunc: "/", ModVal: 0}
	_, err := RollParsed(pd, testRNG)
	if err == nil {
		t.Fatalf("expected division by zero error")
	}
}

func TestExploding(t *testing.T) {
	// Using a seeded RNG that will produce a 6 to cause explosion for a d6
	pr, err := Parse("2d6!")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	res, err := RollParsed(pr, testRNG)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.AllRolls) != 2 {
		t.Fatalf("expected 2 all-rolls, got %d", len(res.AllRolls))
	}
	if len(res.Rolls) != 2 {
		t.Fatalf("expected 2 final rolls, got %d", len(res.Rolls))
	}
}

func TestKeepHighest(t *testing.T) {
	pr, err := Parse("4d6kh3")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	res, err := RollParsed(pr, testRNG)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.AllRolls) != 4 {
		t.Fatalf("expected 4 all-rolls, got %d", len(res.AllRolls))
	}
	if len(res.Rolls) != 3 {
		t.Fatalf("expected 3 kept rolls, got %d", len(res.Rolls))
	}
}

func TestRerollUntil(t *testing.T) {
	pr, err := Parse("4d6r1")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	// use a fresh RNG for determinism
	rng := rand.New(rand.NewPCG(600, 601))
	res, err := RollParsed(pr, rng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// when rerolling until 1, no final per-die total should equal 1
	for _, v := range res.AllRolls {
		if v == 1 {
			t.Fatalf("found value 1 in AllRolls despite reroll rule: %+v", res.AllRolls)
		}
	}
}

func TestRerollOnce(t *testing.T) {
	pr, err := Parse("4d6ro1")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	rng := rand.New(rand.NewPCG(600, 601))
	res, err := RollParsed(pr, rng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.RerollsPerformed < 0 || res.RerollsPerformed > pr.Count {
		t.Fatalf("unexpected rerolls performed: %d", res.RerollsPerformed)
	}
}

func TestSuccessCounting(t *testing.T) {
	pr, err := Parse("10d10>=8")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	rng := rand.New(rand.NewPCG(600, 601))
	res, err := RollParsed(pr, rng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Successes < 0 || res.Successes > pr.Count {
		t.Fatalf("success count out of bounds: %d", res.Successes)
	}
}

func TestParseRerollComparators(t *testing.T) {
	cases := []struct {
		in    string
		op    string
		val   int
		once  bool
		count int
	}{
		{"4d6r<2", "<", 2, false, 0},
		{"3d6r>=5", ">=", 5, false, 0},
		{"2d6ro1#1", "=", 1, true, 1},
		{"4d6r1#2", "=", 1, false, 2},
		{"3d6r!=3", "!=", 3, false, 0},
	}
	for _, c := range cases {
		pd, err := Parse(c.in)
		if err != nil {
			t.Fatalf("parse %q err: %v", c.in, err)
		}
		if pd.RerollOp != c.op || pd.RerollVal != c.val || pd.RerollOnce != c.once || pd.RerollCount != c.count {
			t.Fatalf("Parse(%q) = %+v, want op=%q val=%d once=%v count=%d", c.in, pd, c.op, c.val, c.once, c.count)
		}
	}
}

func TestRerollComparatorBehavior(t *testing.T) {
	pr, err := Parse("4d6r<2")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	rng := rand.New(rand.NewPCG(600, 601))
	res, err := RollParsed(pr, rng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, v := range res.AllRolls {
		if v == 1 {
			t.Fatalf("found value 1 in AllRolls despite r<2: %+v", res.AllRolls)
		}
	}
}

func TestRerollCountLimit(t *testing.T) {
	pr, err := Parse("4d6r1#2")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	if pr.RerollCount != 2 {
		t.Fatalf("expected reroll count 2, got %d", pr.RerollCount)
	}
	rng := rand.New(rand.NewPCG(600, 601))
	res, err := RollParsed(pr, rng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.RerollsPerformed < 0 || res.RerollsPerformed > pr.Count*pr.RerollCount {
		t.Fatalf("unexpected rerolls performed: %d", res.RerollsPerformed)
	}
}

func TestRerollComparatorsVarious(t *testing.T) {
	cases := []struct {
		in        string
		forbidden func(int) bool
	}{
		{"4d6r<=2", func(v int) bool { return v <= 2 }},
		{"4d6r>4", func(v int) bool { return v > 4 }},
		{"4d6ro<2", func(v int) bool { return v < 2 }},
	}
	for _, c := range cases {
		pr, err := Parse(c.in)
		if err != nil {
			t.Fatalf("parse %q err: %v", c.in, err)
		}
		rng := rand.New(rand.NewPCG(600, 601))
		res, err := RollParsed(pr, rng)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", c.in, err)
		}
		for _, v := range res.AllRolls {
			if c.forbidden(v) {
				t.Fatalf("value %d forbidden by %q, rolls: %+v", v, c.in, res.AllRolls)
			}
		}
		// ro should perform at most count rerolls (1 per die)
		if pr.RerollOnce && res.RerollsPerformed > pr.Count {
			t.Fatalf("ro performed too many rerolls: %d > %d", res.RerollsPerformed, pr.Count)
		}
	}
}

func TestRerollInfiniteLoopGuardAndCountStop(t *testing.T) {
	// infinite loop without per-die limit should error (1-sided die rerolling its only face)
	pr, err := Parse("1d1r1")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	if _, err := RollParsed(pr, rand.New(rand.NewPCG(600, 601))); err == nil {
		t.Fatalf("expected error due to exceeding max roll limit for %q", "1d1r1")
	}

	// with per-die limit, it should stop and not error
	pr2, err := Parse("1d1r1#2")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	res, err := RollParsed(pr2, rand.New(rand.NewPCG(600, 601)))
	if err != nil {
		t.Fatalf("unexpected error with per-die limit: %v", err)
	}
	if res.RerollsPerformed != pr2.RerollCount {
		t.Fatalf("expected %d rerolls, got %d", pr2.RerollCount, res.RerollsPerformed)
	}
}

func TestRerollPrecedesExplode(t *testing.T) {
	// r6 should eliminate 6s and thus prevent explosions
	pr, err := Parse("2d6!r6")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	res, err := RollParsed(pr, rand.New(rand.NewPCG(600, 601)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, v := range res.AllRolls {
		if v >= 6 {
			t.Fatalf("expected no value >=6 due to r6 preventing explosion, got %+v", res.AllRolls)
		}
	}
}

func TestRerollWithKeepDrop(t *testing.T) {
	pr, err := Parse("4d6kh3r1")
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	res, err := RollParsed(pr, rand.New(rand.NewPCG(600, 601)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, v := range res.Rolls {
		if v == 1 {
			t.Fatalf("kept roll equals 1 despite r1: %+v", res.Rolls)
		}
	}
}

func TestLimits(t *testing.T) {
	// Too many dice
	large := ParsedDice{Count: MaxDiceCount + 1, Sides: 6, Type: "d"}
	if _, err := RollParsed(large, testRNG); err == nil {
		t.Fatalf("expected error for too many dice")
	}
	// Invalid sides
	bad := ParsedDice{Count: 1, Sides: 0, Type: "d"}
	if _, err := RollParsed(bad, testRNG); err == nil {
		t.Fatalf("expected error for invalid sides")
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Parse("3d6+2")
	}
}

func BenchmarkRollParsed(b *testing.B) {
	pd := ParsedDice{Count: 10, Sides: 8, Type: "d"}
	for i := 0; i < b.N; i++ {
		_, _ = RollParsed(pd, testRNG)
	}
}

func BenchmarkParseExplode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Parse("2d6!")
	}
}

func BenchmarkRollExplode(b *testing.B) {
	pd, _ := Parse("2d6!")
	rng := rand.New(rand.NewPCG(600, 601))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = RollParsed(pd, rng)
	}
}

func BenchmarkRollKeepDrop(b *testing.B) {
	pd, _ := Parse("4d6kh3")
	rng := rand.New(rand.NewPCG(600, 601))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = RollParsed(pd, rng)
	}
}

func BenchmarkRollManyDice(b *testing.B) {
	pd, _ := Parse("100d6")
	rng := rand.New(rand.NewPCG(600, 601))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = RollParsed(pd, rng)
	}
}
