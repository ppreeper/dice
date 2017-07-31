package main

import "testing"

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
	{"1d6", 1, "d", 6, "", 0, 4, []int{6}, 6},
	{"2d6", 2, "d", 6, "", 0, 6, []int{6, 4}, 10},
	{"10d8", 10, "d", 8, "", 0, 7, []int{2, 8, 8, 4, 2, 7, 2, 5, 1, 5}, 44},
	{"4F", 4, "F", 3, "", 0, 3, []int{3, 1, 3, 3}, 2},
	{"d6", 1, "d", 6, "", 0, 1, []int{6}, 6},
	{"d20", 1, "d", 20, "", 0, 12, []int{2}, 2},
	{"d6+2", 1, "d", 6, "+", 2, 6, []int{6}, 8},
	{"d20+5", 1, "d", 20, "+", 5, 15, []int{2}, 7},
	{"3d20x5", 3, "d", 20, "x", 5, 6, []int{2, 8, 8}, 90},
	{"3d20/5", 3, "d", 20, "/", 5, 7, []int{2, 8, 8}, 13},
	{"3d20+5", 3, "d", 20, "+", 5, 7, []int{2, 8, 8}, 23},
	{"3d20-5", 3, "d", 20, "-", 5, 19, []int{2, 8, 8}, 13},
}

// TestPattern test
func TestPattern(t *testing.T) {
	for _, mt := range mulPatternTests {
		var d Dice
		d.seed = 600
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

// TestRollDie test
func TestRollDie(t *testing.T) {
	for _, mt := range mulPatternTests {
		var d Dice
		d.seed = 600
		d.Pattern(mt.a)
		d.RollDie()
		if v := d.RollDie(); mt.expectedRoll != v {
			t.Errorf("\nCount %d, Sides %d, Expected %v, got %v",
				d.DieCount, d.DieSides, mt.expectedRoll, v)
		}
	}
}

// TestRoll test
func TestRoll(t *testing.T) {
	for _, mt := range mulPatternTests {
		var d Dice
		d.seed = 1
		d.Pattern(mt.a)
		d.Roll()
		// fmt.Printf("%v %v\n", d.Results, mt.expected)
		if len(d.Results) == len(mt.expected) {
			for i := 0; i < len(d.Results); i++ {
				if d.Results[i] != mt.expected[i] {
					t.Errorf("results not equal\n")
				}
			}
		}
		if d.Sum != mt.expectedSum {
			t.Errorf("sum not equal\n")
		}
	}
}
