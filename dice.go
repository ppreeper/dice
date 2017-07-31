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

var randSource = rand.NewSource(time.Now().UnixNano())
var randFixed = rand.NewSource(600)

// Dice structure of the dice set
type Dice struct {
	DieCount   int
	DieType    string
	DieSides   int
	DieModFunc string
	DieModVal  int
	Results    []int
	total      int
	seed       bool
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
func (d *Dice) RollDie(r *rand.Rand) int {
	var rollVal int
	rollVal = r.Intn(d.DieSides) + 1
	return rollVal
}

// Sum channel
func (d *Dice) Sum(s []int, c chan int) {
	sum := 0
	// fmt.Printf("%v\n", s)
	for _, v := range s {
		// fmt.Printf("%d %d\n", k, v)
		sum += v
	}
	// fmt.Printf("in chan %d\n", sum)
	c <- sum
}

// Roll soll the die given
func (d *Dice) Roll(die string) {
	d.Pattern(die)
	c := make(chan int)
	// fmt.Printf("DieCount = %d, DieSides = %d\n", d.DieCount, d.DieSides)
	var r *rand.Rand
	if d.seed {
		r = rand.New(randFixed)
	} else {
		r = rand.New(randSource)
	}
	for i := 0; i < d.DieCount; i++ {
		// d.RollDie(r)
		res := d.RollDie(r)
		d.Results = append(d.Results, res)
		// fmt.Printf("DieSides = %d, result=%d\n", d.DieSides, r)
	}
	// fmt.Println(d.Results)
	if d.DieType == "F" {
		for i := 0; i < len(d.Results); i++ {
			switch d.Results[i] {
			case 1:
				d.Results[i] = -1
			case 2:
				d.Results[i] = 0
			case 3:
				d.Results[i] = 1
			}
		}
	}

	// if d.seed == false {
	// 	randDie.NewSource(time.Now().UnixNano())
	go d.Sum(d.Results, c)
	d.total = <-c
	// } else {
	// 	go d.Sum(d.Results, c)
	// 	d.total = <-c
	// }

	// fmt.Println(d.total)
	switch d.DieModFunc {
	case "+":
		// fmt.Printf("%v %d add %d = %d\n", d.Results, d.Sum, d.DieModVal, d.Sum+d.DieModVal)
		d.total = d.total + d.DieModVal
	case "-":
		// fmt.Printf("%v %d minus %d = %d\n", d.Results, d.Sum, d.DieModVal, d.Sum-d.DieModVal)
		d.total = d.total - d.DieModVal
	case "x":
		// fmt.Printf("%v %d times %d = %d\n", d.Results, d.Sum, d.DieModVal, d.Sum*d.DieModVal)
		d.total = d.total * d.DieModVal
	case "/":
		// fmt.Printf("%v %d divide %d = %d\n", d.Results, d.Sum, d.DieModVal, d.Sum/d.DieModVal)
		d.total = d.total / d.DieModVal
	}
}

func main() {
	flag.Parse()
	var d Dice
	d.seed = false
	if *cdn == "" {
		fmt.Println("No Die Notation")
	} else {
		d.Roll(*cdn)
		fmt.Printf("%v %v\n", d.Results, d.total)
	}
	return
}
