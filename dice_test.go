package dice

import (
	"math/rand/v2"
	"testing"
)

// small, deterministic RNG for tests
var testRNG = rand.New(rand.NewPCG(600, 601))

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
