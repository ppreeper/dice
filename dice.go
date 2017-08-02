package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
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

// Pattern2 determine pattern from dice notation string refactor
func (d *Dice) Pattern2(die string) {
	switch {
	case strings.HasSuffix(die, "F"):
		d.DieType = "F"
		d.DieCount, _ = strconv.Atoi(strings.TrimSuffix(die, "F"))
		d.DieSides = 3
		d.DieModFunc = ""
		d.DieModVal = 0
	case strings.HasPrefix(die, "d"):
		d.DieType = "d"
		d.DieCount = 1
		remainder := strings.TrimPrefix(die, "d")
		switch {
		case strings.Contains(remainder, "+"):
			remVals := strings.Split(remainder, "+")
			d.DieSides, _ = strconv.Atoi(remVals[0])
			d.DieModFunc = "+"
			d.DieModVal, _ = strconv.Atoi(remVals[1])
		case strings.Contains(remainder, "-"):
			remVals := strings.Split(remainder, "-")
			d.DieSides, _ = strconv.Atoi(remVals[0])
			d.DieModFunc = "-"
			d.DieModVal, _ = strconv.Atoi(remVals[1])
		case strings.Contains(remainder, "x"):
			remVals := strings.Split(remainder, "x")
			d.DieSides, _ = strconv.Atoi(remVals[0])
			d.DieModFunc = "x"
			d.DieModVal, _ = strconv.Atoi(remVals[1])
		case strings.Contains(remainder, "/"):
			remVals := strings.Split(remainder, "/")
			d.DieSides, _ = strconv.Atoi(remVals[0])
			d.DieModFunc = "/"
			d.DieModVal, _ = strconv.Atoi(remVals[1])
		default:
			d.DieSides, _ = strconv.Atoi(remainder)
			d.DieModFunc = ""
			d.DieModVal = 0
		}
	case strings.Contains(die, "d"):
		d.DieType = "d"
		dieSplit := strings.Split(die, "d")
		d.DieCount, _ = strconv.Atoi(dieSplit[0])
		remainder := dieSplit[1]
		switch {
		case strings.Contains(remainder, "+"):
			remVals := strings.Split(remainder, "+")
			d.DieSides, _ = strconv.Atoi(remVals[0])
			d.DieModFunc = "+"
			d.DieModVal, _ = strconv.Atoi(remVals[1])
		case strings.Contains(remainder, "-"):
			remVals := strings.Split(remainder, "-")
			d.DieSides, _ = strconv.Atoi(remVals[0])
			d.DieModFunc = "-"
			d.DieModVal, _ = strconv.Atoi(remVals[1])
		case strings.Contains(remainder, "x"):
			remVals := strings.Split(remainder, "x")
			d.DieSides, _ = strconv.Atoi(remVals[0])
			d.DieModFunc = "x"
			d.DieModVal, _ = strconv.Atoi(remVals[1])
		case strings.Contains(remainder, "/"):
			remVals := strings.Split(remainder, "/")
			d.DieSides, _ = strconv.Atoi(remVals[0])
			d.DieModFunc = "/"
			d.DieModVal, _ = strconv.Atoi(remVals[1])
		default:
			d.DieSides, _ = strconv.Atoi(remainder)
			d.DieModFunc = ""
			d.DieModVal = 0
		}
	}
}

// Pattern determine pattern from dice notation string
func (d *Dice) Pattern(die string) {
	dCount := regexp.MustCompile(`(\d+)d`)
	sCount := regexp.MustCompile(`d(\d+)`)
	modifier := regexp.MustCompile(`[x/+-]`)
	modVal := regexp.MustCompile(`[x/+-](\d+)`)
	sPattern := regexp.MustCompile(`^(\d+)?d(\d+)([x/+-](\d+))?`)
	fudge := regexp.MustCompile(`^(\d+)F`)

	if sPattern.MatchString(die) == true {
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
	for _, v := range s {
		sum += v
	}
	c <- sum
}

// Roll soll the die given
func (d *Dice) Roll(die string) {
	d.Pattern(die)
	c := make(chan int)
	var r *rand.Rand
	if d.seed {
		r = rand.New(randFixed)
	} else {
		r = rand.New(randSource)
	}
	for i := 0; i < d.DieCount; i++ {
		res := d.RollDie(r)
		d.Results = append(d.Results, res)
	}
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

	go d.Sum(d.Results, c)
	d.total = <-c

	switch d.DieModFunc {
	case "+":
		d.total = d.total + d.DieModVal
	case "-":
		d.total = d.total - d.DieModVal
	case "x":
		d.total = d.total * d.DieModVal
	case "/":
		d.total = d.total / d.DieModVal
	}
}

func main() {
	flag.Parse()
	var d Dice
	d.seed = false
	if *cdn == "" {
		fmt.Fprintf(os.Stdout, "No Die Notation\n")
	} else {
		d.Roll(*cdn)
		fmt.Fprintf(os.Stdout, "%v %v\n", d.Results, d.total)
	}
	return
}
