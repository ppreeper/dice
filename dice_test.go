package dice

import (
	"bytes"
	"io"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

// patternTest structure of the dice pattern tests
type patternTest struct {
	a          string
	dieCount   int
	dieType    string
	dieSides   int
	dieModFunc string
	dieModVal  int
}

func dicePatGen() []patternTest {
	// var dice []string
	var patTests []patternTest
	for i := 0; i < 101; i++ {
		d := strconv.Itoa(i) + "F"
		// dice = append(dice, d)
		patTests = append(patTests, patternTest{d, i, "F", 3, "", 0})
	}
	for i := 0; i < 101; i++ {
		d := "d" + strconv.Itoa(i)
		// dice = append(dice, d)
		patTests = append(patTests, patternTest{d, 1, "d", i, "", 0})
	}
	for i := 0; i < 101; i++ {
		for j := 0; j < 101; j++ {
			d := strconv.Itoa(i) + "d" + strconv.Itoa(j)
			// dice = append(dice, d)
			patTests = append(patTests, patternTest{d, i, "d", j, "", 0})
		}
	}
	for i := 0; i < 101; i++ {
		for _, m := range []string{"+", "-", "x", "/"} {
			for mv := 0; mv < 101; mv++ {
				d := "d" + strconv.Itoa(i) + m + strconv.Itoa(mv)
				// dice = append(dice, d)
				patTests = append(patTests, patternTest{d, 1, "d", i, m, mv})
			}
		}
	}
	for i := 0; i < 101; i++ {
		for j := 0; j < 101; j++ {
			for _, m := range []string{"+", "-", "x", "/"} {
				for mv := 0; mv < 101; mv++ {
					d := strconv.Itoa(i) + "d" + strconv.Itoa(j) + m + strconv.Itoa(mv)
					// dice = append(dice, d)
					patTests = append(patTests, patternTest{d, i, "d", j, m, mv})
				}
			}
		}
	}
	return patTests
}

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
	{"1d6", 1, "d", 6, "", 0, 6, []int{5}, 6},
	{"2d6", 2, "d", 6, "", 0, 6, []int{3, 4}, 7},
	{"10d8", 10, "d", 8, "", 0, 6, []int{2, 5, 4, 5, 7, 8, 1, 6, 5, 2}, 44},
	{"4F", 4, "F", 3, "", 0, 2, []int{-1, 0, 0, 1}, 0},
	{"d6", 1, "d", 6, "", 0, 3, []int{1}, 1},
	{"d20", 1, "d", 20, "", 0, 16, []int{8}, 8},
	{"d6+2", 1, "d", 6, "+", 2, 5, []int{5}, 7},
	{"d6-2", 1, "d", 6, "-", 2, 1, []int{4}, 2},
	{"d6x2", 1, "d", 6, "x", 2, 4, []int{3}, 4},
	{"d6/2", 1, "d", 6, "/", 2, 1, []int{2}, 1},
	{"d20+5", 1, "d", 20, "+", 5, 16, []int{9}, 16},
	{"3d20x5", 3, "d", 20, "x", 5, 3, []int{8, 9, 20}, 250},
	{"3d20/5", 3, "d", 20, "/", 5, 9, []int{4, 5, 11}, 2},
	{"3d20+5", 3, "d", 20, "+", 5, 3, []int{15, 18, 17}, 46},
	{"3d20-5", 3, "d", 20, "-", 5, 11, []int{3, 2, 8}, 30},
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
			t.Errorf("\nDieModFunc %s expected %s, got %s", mt.a, mt.dieModFunc, d.DieModFunc)
		}
		if d.DieModVal != mt.dieModVal {
			t.Errorf("\nDieModVal expected %d, got %d", mt.dieModVal, d.DieModVal)
		}
	}
}

// TestPattern2 test
func TestPattern2(t *testing.T) {
	for _, mt := range dicePatGen() {
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
			t.Errorf("\nDieModFunc %s expected %s, got %s", mt.a, mt.dieModFunc, d.DieModFunc)
		}
		if d.DieModVal != mt.dieModVal {
			t.Errorf("\nDieModVal expected %d, got %d", mt.dieModVal, d.DieModVal)
		}
	}
}

func BenchmarkPattern_6F(b *testing.B) {
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("6F")
	}
}

func BenchmarkPattern_d6(b *testing.B) {
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("d6")
	}
}

func BenchmarkPattern_1d6(b *testing.B) {
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("1d6")
	}
}

func BenchmarkPattern_d6add(b *testing.B) {
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("d6+1")
	}
}
func BenchmarkPattern_d6sub(b *testing.B) {
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("d6-1")
	}
}
func BenchmarkPattern_d6mul(b *testing.B) {
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("d6x1")
	}
}
func BenchmarkPattern_d6div(b *testing.B) {
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("d6/1")
	}
}
func BenchmarkPattern_2d6add(b *testing.B) {
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("2d6+1")
	}
}
func BenchmarkPattern_2d6sub(b *testing.B) {
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("2d6-1")
	}
}
func BenchmarkPattern_2d6mul(b *testing.B) {
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("2d6x1")
	}
}
func BenchmarkPattern_2d6div(b *testing.B) {
	var d Dice
	for n := 0; n < b.N; n++ {
		d.Pattern("2d6/1")
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

// BenchmarkRollDie testing
func BenchmarkRollDie(b *testing.B) {
	cdn := "1d6"
	var d Dice
	var r *rand.Rand
	r = rand.New(randFixed)
	d.seed = true
	d.Pattern(cdn)
	for n := 0; n < b.N; n++ {
		d.RollDie(r)
	}
}

// TestRoll test
func TestRoll(t *testing.T) {
	for _, mt := range mulPatternTests {
		var d Dice
		d.seed = true
		d.Roll(mt.a)
		if len(d.Results) == len(mt.expected) {
			if d.total != mt.expectedSum {
				t.Errorf("results %v\t%v\tsum not equal %d  %d\n", d.Results, mt.expected, d.total, mt.expectedSum)
			}
		}
	}
}

// BenchmarkRoll testing
func BenchmarkRoll(b *testing.B) {
	cdn := "1d6"
	var d Dice
	d.seed = true
	for n := 0; n < b.N; n++ {
		d.Roll(cdn)
	}
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
