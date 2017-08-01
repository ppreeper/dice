package main

import (
	"math/rand"
	"testing"
)

// mulPatternTests use values that you know are right
var mulPatternTests = []struct {
	a            string
	dieCount     int
	dieType      string
	dieSides     int
	dieModFunc   string
	dieModVal    int
	expectedRoll int
	expected     []int
	expectedSum  int
}{
	{"1d6", 1, "d", 6, "", 0, 6, []int{5}, 5},
	{"2d6", 2, "d", 6, "", 0, 6, []int{3, 3}, 6},
	{"10d8", 10, "d", 8, "", 0, 6, []int{2, 5, 4, 5, 7, 8, 1, 6, 5, 2}, 45},
	{"4F", 4, "F", 3, "", 0, 2, []int{0, 1, -1, -1}, -1},
	{"d6", 1, "d", 6, "", 0, 3, []int{2}, 2},
	{"d20", 1, "d", 20, "", 0, 16, []int{11}, 11},
	{"d6+2", 1, "d", 6, "+", 2, 5, []int{3}, 5},
	{"d20+5", 1, "d", 20, "+", 5, 15, []int{9}, 14},
	{"3d20x5", 3, "d", 20, "x", 5, 12, []int{8, 9, 20}, 185},
	{"3d20/5", 3, "d", 20, "/", 5, 1, []int{4, 5, 11}, 4},
	{"3d20+5", 3, "d", 20, "+", 5, 16, []int{15, 18, 17}, 55},
	{"3d20-5", 3, "d", 20, "-", 5, 3, []int{3, 2, 8}, 8},
}

// TestPattern test
func TestPattern(t *testing.T) {
	for _, mt := range mulPatternTests {
		var d Dice
		d.seed = true
		d.Pattern(mt.a)
		if d.DieType != mt.dieType {
			t.Errorf("\nDieType expected %s, got %s", mt.dieType, d.DieType)
		}
		if d.DieCount != mt.dieCount {
			t.Errorf("\nDieCount expected %d, got %d", mt.dieCount, d.DieCount)
		}
		if d.DieSides != mt.dieSides {
			t.Errorf("\nDieSides expected %d, got %d", mt.dieSides, d.DieSides)
		}
		if d.DieModFunc != mt.dieModFunc {
			t.Errorf("\nDieModFunc expected %s, got %s", mt.dieModFunc, d.DieModFunc)
		}
		if d.DieModVal != mt.dieModVal {
			t.Errorf("\nDieModVal expected %d, got %d", mt.dieModVal, d.DieModVal)
		}
	}
}

// BenchmarkPattern
func BenchmarkPatternd6(b *testing.B) {
	// run the Fib function b.N times
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("d6")
	}
}

// BenchmarkPattern
func BenchmarkPattern1d6(b *testing.B) {
	// run the Fib function b.N times
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("1d6")
	}
}

func BenchmarkPattern6F(b *testing.B) {
	// run the Fib function b.N times
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("6F")
	}
}

// TestRollDie test
func TestRollDie(t *testing.T) {
	for _, mt := range mulPatternTests {
		var d Dice
		d.seed = true
		var r *rand.Rand
		r = rand.New(randFixed)
		d.Pattern(mt.a)
		// v := d.RollDie(r)
		if v := d.RollDie(r); mt.expectedRoll != v {
			t.Errorf("\nCount %d, Sides %d, Expected %v, got %v",
				d.DieCount, d.DieSides, mt.expectedRoll, v)
		}
	}
}

// TestRoll test
func TestRoll(t *testing.T) {
	for _, mt := range mulPatternTests {
		var d Dice
		d.seed = true
		d.Roll(mt.a)
		// fmt.Printf("%v %v\n", d.Results, mt.expected)
		// fmt.Printf("result lengths %d  %d\t", len(d.Results), len(mt.expected))
		// fmt.Printf("sum %d  %d\n", d.total, mt.expectedSum)
		if len(d.Results) == len(mt.expected) {
			if d.total != mt.expectedSum {
				t.Errorf("results %v\t%v\tsum not equal %d  %d\n", d.Results, mt.expected, d.total, mt.expectedSum)
			}
		}
	}
}

// TestMain testing
func TestMain(t *testing.T) {
	*cdn = "1d6"
	main()
	*cdn = ""
	main()
}
