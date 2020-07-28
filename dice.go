package dice

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

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

// Pattern determine pattern from dice notation string refactor
func (d *Dice) Pattern(die string) {
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

// RollDie soll the die given
func (d *Dice) RollDie(r *rand.Rand) int {
	var rollVal int
	rollVal = r.Intn(d.DieSides) + 1
	return rollVal
}

// Roll the die given
func (d *Dice) Roll(die string) {
	d.Pattern(die)
	var r *rand.Rand
	if d.seed {
		r = rand.New(randFixed)
	} else {
		r = rand.New(randSource)
	}
	for i := 0; i < d.DieCount; i++ {
		res := d.RollDie(r)
		if d.DieType == "F" {
			d.total += res - 2
		} else {
			d.total += res
		}
		d.Results = append(d.Results, res)
	}

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
