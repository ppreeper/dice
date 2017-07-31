package main

import (
	"flag"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

var cdn = flag.String("D", "", "Die Notation")

// Dice structure of the dice set
type Dice struct {
	DieCount   int
	DieType    string
	DieSides   int
	DieModFunc string
	DieModVal  int
	Results    []int
	Sum        int
	seed       int64
}

// Pattern determine pattern from dice notation string
func (d *Dice) Pattern(die string) {
	dCount := regexp.MustCompile(`(\d+)d`)
	sCount := regexp.MustCompile(`d(\d+)`)
	modifier := regexp.MustCompile(`[x/+-]`)
	modVal := regexp.MustCompile(`[x/+-](\d+)`)
	sPattern := regexp.MustCompile(`^(\d+)?d(\d+)([x/+-](\d+))?`)
	fudge := regexp.MustCompile(`^(\d+)F`)

	// fmt.Printf("\n%s\t", die)
	if sPattern.MatchString(die) == true {
		// fmt.Printf("sPattern\t")
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
		d.DieSides, _ = strconv.Atoi(string(sc[1]))
	}
	if fudge.MatchString(die) == true {
		d.DieType = "F"
		dc := fudge.FindSubmatch([]byte(die))
		d.DieCount, _ = strconv.Atoi(string(dc[1]))
		d.DieSides = 3
		d.DieModFunc = ""
		d.DieModVal = 0
	}
}

// RollDie soll the die given
func (d *Dice) RollDie() int {
	rand.New(rand.NewSource(d.seed))
	rollVal := rand.Intn(d.DieSides) + 1
	// fmt.Printf("%v %v %v\n", d.DieCount, d.DieSides, rollVal)
	return rollVal
}

// Roll soll the die given
func (d *Dice) Roll() {
	rand.Seed(d.seed)
	// fmt.Printf("DieCount = %d, DieSides = %d\n", d.DieCount, d.DieSides)
	for i := 0; i < d.DieCount; i++ {
		r := d.RollDie()
		d.Results = append(d.Results, r)
		// fmt.Printf("DieSides = %d, result=%d\n", d.DieSides, r)
	}
	// fmt.Println(d.Results)
	if d.DieType == "F" {
		d.Sum = 0
		for i := 0; i < len(d.Results); i++ {
			switch d.Results[i] {
			case 1:
				d.Sum = d.Sum - 1
			case 3:
				d.Sum = d.Sum + 1
			}
		}
	} else {
		for i := 0; i < len(d.Results); i++ {
			d.Sum += d.Results[i]
		}
	}
	// fmt.Println(d.Sum)
	switch d.DieModFunc {
	case "+":
		// fmt.Printf("%v %d add %d = %d\n", d.Results, d.Sum, d.DieModVal, d.Sum+d.DieModVal)
		d.Sum = d.Sum + d.DieModVal
	case "-":
		// fmt.Printf("%v %d minus %d = %d\n", d.Results, d.Sum, d.DieModVal, d.Sum-d.DieModVal)
		d.Sum = d.Sum - d.DieModVal
	case "x":
		// fmt.Printf("%v %d times %d = %d\n", d.Results, d.Sum, d.DieModVal, d.Sum*d.DieModVal)
		d.Sum = d.Sum * d.DieModVal
	case "/":
		// fmt.Printf("%v %d divide %d = %d\n", d.Results, d.Sum, d.DieModVal, d.Sum/d.DieModVal)
		d.Sum = d.Sum - d.DieModVal
	}
}

func main() {
	flag.Parse()
	var d Dice
	d.seed = time.Now().UnixNano()
	if *cdn == "" {
		fmt.Println("No Die Notation")
	} else {
		d.Pattern(*cdn)
		d.Roll()
		fmt.Printf("%v %v\n", d.Results, d.Sum)
	}
	return
}
