package main

import (
	"testing"
)

var mulRoll = []struct {
	a        int
	expected int
}{
	{1, 0},
	{2, 1},
	{3, 2},
	{4, 3},
	{5, 1},
	{6, 5},
	{7, 0},
	{8, 3},
	{9, 5},
	{10, 1},
	{11, 6},
	{12, 11},
	{13, 10},
	{14, 7},
	{15, 11},
	{16, 3},
	{17, 15},
	{18, 5},
	{19, 14},
	{20, 11},
}

func TestRoll1(t *testing.T) {
	for _, mt := range mulRoll {
		if v := roll1(mt.a); v != mt.expected {
			t.Errorf("roll1(%d) returned %d, expected %d", mt.a, v, mt.expected)
		}
	}
}

// use values that you know are right
var mulPatternTests = []struct {
	a          string
	dieCount   int
	dieType    string
	dieSides   int
	dieModFunc string
	dieModVal  int
	expected   int
}{
	{"1d6", 1, "d", 6, "", 0, 5},
	{"2d6", 2, "d", 6, "", 0, 5},
	{"10d8", 10, "d", 8, "", 0, 3},
	{"4F", 4, "F", 3, "", 0, 2},
	{"d6", 1, "d", 6, "", 0, 5},
	{"d20", 1, "d", 20, "", 0, 11},
	{"d6+2", 1, "d", 6, "+", 2, 5},
	{"d20+5", 1, "d", 20, "+", 5, 11},
	{"3d20x5", 3, "d", 20, "x", 5, 11},
	{"3d20/5", 3, "d", 20, "/", 5, 11},
	{"3d20+5", 3, "d", 20, "+", 5, 11},
	{"3d20-5", 3, "d", 20, "-", 5, 11},
}

// type Dice struct {
// 	DieCount   int
// 	DieType    string
// 	DieSides   int
// 	DieModFunc string
// 	DieModVal  int
// }

func TestPattern(t *testing.T) {
	var d Dice
	for _, mt := range mulPatternTests {
		// fmt.Printf("%s", mt.a)
		d.Pattern(mt.a)
		// fmt.Printf("\tDieType: %s", d.DieType)
		// fmt.Printf("\tDieCount: %d", d.DieCount)
		// fmt.Printf("\tDieSides: %d", d.DieSides)
		// fmt.Printf("\tDieModFunc: %s", d.DieModFunc)
		// fmt.Printf("\tDieModVal: %d\n", d.DieModVal)
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

func TestRoll(t *testing.T) {
	var d Dice
	for _, mt := range mulPatternTests {
		d.Pattern(mt.a)
		if v := d.Roll(d.DieSides); v != mt.expected {
			t.Errorf("\nDie.Roll(%d) expected %d, got %d",
				d.DieSides,
				mt.expected,
				d.DieModVal)
		}
	}
}
