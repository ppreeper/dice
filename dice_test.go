package main

import (
	"fmt"
	"testing"
)

// use values that you know are right
var mulTests = []struct {
	a, b     int
	expected int
}{
	{1, 1, 1},
	{2, 2, 4},
	{3, 3, 9},
	{4, 4, 16},
	{5, 5, 25},
}

func TestMul(t *testing.T) {
	for _, mt := range mulTests {
		if v := Mul(mt.a, mt.b); v != mt.expected {
			t.Errorf("Mul(%d, %d) returned %d, expected %d", mt.a, mt.b, v, mt.expected)
		}
	}
}

// use values that you know are right
var fibTests = []uint64{1, 1, 2, 3, 5, 8, 13, 21, 34, 55}

func TestFibFunc(t *testing.T) {
	fn := FibFunc()
	for i, v := range fibTests {
		if val := fn(); val != v {
			t.Fatalf("at index %d, expected %d, got %d.", i, v, val)
		}
	}
}

func BenchmarkFibFunc(b *testing.B) {
	fn := FibFunc()
	for i := 0; i < b.N; i++ {
		_ = fn()
	}
}

var di = []string{
	"1d6", "2d6", "10d8",
	"4F",
	"d6", "d20",
	"d6+2", "d20+5",
	"3d20x5", "3d20/5",
	"3d20+5", "3d20-5",
}

// use values that you know are right
// var mulPatternTests = []struct {
// 	a        string
// 	expected Dice
// }{
// 	{"1d6", 1},
// 	{"2d6", 1},
// 	{"10d8", 1},
// 	{"4F", 1},
// 	{"d6", 1},
// 	{"d20", 1},
// 	{"d6+2", 1},
// 	{"d20+5", 1},
// 	{"3d20x5", 1},
// 	{"3d20/5", 1},
// 	{"3d20+5", 1},
// 	{"3d20-5", 1},
// }

func TestPattern(t *testing.T) {
	var d Dice
	d.Pattern("d6")
	fmt.Printf("DieType: %s", d.DieType)
	fmt.Printf("\tDieCount: %d", d.DieCount)
	if d.DieType != "d" {
		t.Errorf("Pattern expected d")
	}
	d.Pattern("3d6")
	fmt.Printf("DieType: %s", d.DieType)
	fmt.Printf("\tDieCount: %d", d.DieCount)
	d.Pattern("d6+20")
	fmt.Printf("DieType: %s", d.DieType)
	fmt.Printf("\tDieCount: %d", d.DieCount)
	d.Pattern("3d6+20")
	fmt.Printf("DieType: %s", d.DieType)
	fmt.Printf("\tDieCount: %d", d.DieCount)
	d.Pattern("4F")
	fmt.Printf("DieType: %s", d.DieType)
	fmt.Printf("\tDieCount: %d", d.DieCount)
}
