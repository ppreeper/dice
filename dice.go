package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
)

// Dice structure of the dice set
type Dice struct {
	DieType    string
	DieModFunc string
	DieCount   int
	DieSides   string
	DieModVal  int
}

// Pattern determine pattern from dice notation string
func (d *Dice) Pattern(die string) {
	dCount := regexp.MustCompile(`(\d+)d`)
	sCount := regexp.MustCompile(`d(\d+)`)
	modifier := regexp.MustCompile(`[x/+-]`)
	modVal := regexp.MustCompile(`[x/+-](\d+)`)
	sPattern := regexp.MustCompile(`^(\d+)?d(\d+)([x/+-](\d+))?`)
	fudge := regexp.MustCompile(`^(\d+)F`)

	fmt.Printf("\n%s\t", die)
	if sPattern.MatchString(die) == true {
		fmt.Printf("sPattern\t")
		d.DieType = "d"
		dc := dCount.FindSubmatch([]byte(die))
		sc := sCount.FindSubmatch([]byte(die))
		mdf := modifier.FindSubmatch([]byte(die))
		mdv := modVal.FindSubmatch([]byte(die))
		if len(dc) == 0 {
			d.DieCount = 1
		} else {
			d.DieCount, _ = strconv.Atoi(string(dc[1]))
		}
		if len(mdf) == 0 {
			d.DieModFunc = ""
		} else {
			d.DieModFunc = string(mdf[0])
		}
		if len(mdv) == 0 {
			d.DieModVal = 0
		} else {
			d.DieModVal, _ = strconv.Atoi(string(mdv[1]))
		}
		d.DieSides = string(sc[1])
	}
	if fudge.MatchString(die) == true {
		d.DieType = "F"
		dc := fudge.FindSubmatch([]byte(die))
		d.DieCount, _ = strconv.Atoi(string(dc[1]))
	}
}

//Roll soll the die given
func (d *Dice) Roll(sides int) int {
	// rand.Seed(0)
	return rand.Intn(sides)
}

func main() {
	fmt.Println("Dice")
	// fmt.Println(cdn("1d6"))
	// fmt.Println(cdn("1d6"))
}

//Mul multiply
func Mul(a int, b int) int {
	return a * b
}

//FibFunc fibonacci series
func FibFunc() func() uint64 {
	var a, b uint64 = 0, 1 // yes, it's wrong
	return func() uint64 {
		a, b = b, a+b
		return a
	}
}
