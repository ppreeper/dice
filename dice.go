package dice

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// MaxDiceCount limits the number of dice parsed to avoid excessive allocations.
	MaxDiceCount = 1000
	// MaxSides limits faces of a die to a reasonable upper bound.
	MaxSides = 1000000
	// MaxTotalRolls limits the total number of RNG calls (including exploding)
	// to avoid infinite loops or very large workloads.
	MaxTotalRolls = 10000
)

// ParsedDice represents a parsed dice notation in a minimal, pure form.
type ParsedDice struct {
	Count int // number of dice
	Sides int // number of sides (for Fate dice Sides==3)
	Type  string

	// Explode indicates '!' after sides (basic exploding dice)
	Explode bool

	// Keep/Drop: Action is "" | "k" | "d". Which is "h" or "l" (highest/lowest).
	KeepDropAction string
	KeepDropWhich  string
	KeepDropCount  int

	// Modifiers applied to the final total
	ModFunc string // "", "+", "-", "x", "/"
	ModVal  int

	// Reroll: rN (reroll until not N) or roN (reroll once)
	RerollVal   int
	RerollOnce  bool
	RerollOp    string // comparator: "=", "<", ">", "<=", ">=", "!="
	RerollCount int    // max rerolls per die (0 = unlimited)

	// Success counting: operator and threshold (e.g. ">=8")
	SuccessOp  string
	SuccessVal int
}

// RollResult is a pure result of rolling a ParsedDice with a given RNG.
// AllRolls contains the per-die totals (after exploding). Rolls contains the
// final contributing die values after keep/drop is applied (or equal to AllRolls
// if no keep/drop was requested).
type RollResult struct {
	AllRolls []int
	Rolls    []int
	Total    int
	// Successes counts how many kept rolls meet the success condition (if any)
	Successes int
	// RerollsPerformed counts how many extra rolls were executed due to rerolls
	RerollsPerformed int
	// TotalRollCalls is the total number of RNG calls performed during the roll
	TotalRollCalls int
}

// defaultRand is safe to use concurrently (rand.Rand has internal locking)
// Use PCG seeded from the current time.
var defaultRand = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()>>1)))

// Parse parses a dice notation string into a ParsedDice.
// Supported (expanded) forms include basic exploding (!) and keep/drop (kh/kl/dh/dl or k/d).
// Examples: "4d6kh3", "2d6!", "d20+5", "4F", "d%".
func Parse(s string) (ParsedDice, error) {
	var pd ParsedDice
	s = strings.TrimSpace(s)
	if s == "" {
		return pd, fmt.Errorf("empty notation")
	}

	// Fate dice shorthand: e.g. 4F
	if strings.HasSuffix(s, "F") {
		pref := strings.TrimSuffix(s, "F")
		if pref == "" {
			return pd, fmt.Errorf("missing count before F")
		}
		n, err := strconv.Atoi(pref)
		if err != nil {
			return pd, fmt.Errorf("invalid fate count %q: %w", pref, err)
		}
		if n < 0 || n > MaxDiceCount {
			return pd, fmt.Errorf("count out of range")
		}
		pd.Type = "F"
		pd.Count = n
		pd.Sides = 3 // internal representation for generating -1..1 via 1..3
		return pd, nil
	}

	// Must contain a 'd' for standard dice forms (e.g. "d6", "2d6+1")
	idx := strings.Index(s, "d")
	if idx == -1 {
		return pd, fmt.Errorf("invalid notation (missing 'd' or 'F'): %q", s)
	}

	left := s[:idx]
	right := s[idx+1:]

	// count
	if left == "" {
		pd.Count = 1
	} else {
		n, err := strconv.Atoi(left)
		if err != nil {
			return pd, fmt.Errorf("invalid count %q: %w", left, err)
		}
		pd.Count = n
	}
	if pd.Count < 0 || pd.Count > MaxDiceCount {
		return pd, fmt.Errorf("count out of range: %d", pd.Count)
	}

	// parse sides (digits or '%')
	pos := 0
	if pos >= len(right) {
		return pd, fmt.Errorf("missing sides in notation %q", s)
	}
	if right[pos] == '%' {
		pd.Sides = 100
		pos++
	} else {
		// collect digits
		start := pos
		for pos < len(right) && right[pos] >= '0' && right[pos] <= '9' {
			pos++
		}
		if start == pos {
			return pd, fmt.Errorf("missing sides in notation %q", s)
		}
		sv, err := strconv.Atoi(right[start:pos])
		if err != nil {
			return pd, fmt.Errorf("invalid sides: %w", err)
		}
		pd.Sides = sv
	}
	if pd.Sides <= 0 || pd.Sides > MaxSides {
		return pd, fmt.Errorf("sides out of range: %d", pd.Sides)
	}

	pd.Type = "d"

	// optional explode marker '!'
	if pos < len(right) && right[pos] == '!' {
		pd.Explode = true
		pos++
		// exploding on a die with 1 face is not allowed (would infinite loop)
		if pd.Sides <= 1 {
			return pd, fmt.Errorf("cannot explode on sides <= 1")
		}
	}

	// flexible parsing: accept tokens in any order after optional '!'
	for pos < len(right) {
		ch := right[pos]
		switch ch {
		case 'k', 'd':
			// keep/drop
			pd.KeepDropAction = string(ch)
			pos++
			if pos < len(right) && (right[pos] == 'h' || right[pos] == 'l') {
				pd.KeepDropWhich = string(right[pos])
				pos++
			} else {
				if pd.KeepDropAction == "k" {
					pd.KeepDropWhich = "h"
				} else {
					pd.KeepDropWhich = "l"
				}
			}
			if pos >= len(right) || right[pos] < '0' || right[pos] > '9' {
				return pd, fmt.Errorf("parse error: missing keep/drop count in %q", s)
			}
			start := pos
			for pos < len(right) && right[pos] >= '0' && right[pos] <= '9' {
				pos++
			}
			kv, err := strconv.Atoi(right[start:pos])
			if err != nil {
				return pd, fmt.Errorf("invalid keep/drop count: %w", err)
			}
			pd.KeepDropCount = kv
			continue
		case 'r':
			// reroll
			pos++
			once := false
			if pos < len(right) && right[pos] == 'o' {
				once = true
				pos++
			}
			rerOp := ""
			if pos+1 < len(right) {
				two := right[pos : pos+2]
				if two == ">=" || two == "<=" || two == "!=" {
					rerOp = two
					pos += 2
				}
			}
			if rerOp == "" && pos < len(right) {
				if right[pos] == '>' || right[pos] == '<' || right[pos] == '=' {
					rerOp = string(right[pos])
					pos++
				}
			}
			if rerOp == "" {
				rerOp = "="
			}
			if pos >= len(right) || right[pos] < '0' || right[pos] > '9' {
				return pd, fmt.Errorf("parse error: missing reroll value in %q", s)
			}
			start := pos
			for pos < len(right) && right[pos] >= '0' && right[pos] <= '9' {
				pos++
			}
			rv, err := strconv.Atoi(right[start:pos])
			if err != nil {
				return pd, fmt.Errorf("invalid reroll value: %w", err)
			}
			pd.RerollVal = rv
			pd.RerollOnce = once
			pd.RerollOp = rerOp
			if pos < len(right) && right[pos] == '#' {
				pos++
				if pos >= len(right) || right[pos] < '0' || right[pos] > '9' {
					return pd, fmt.Errorf("parse error: missing reroll count in %q", s)
				}
				start := pos
				for pos < len(right) && right[pos] >= '0' && right[pos] <= '9' {
					pos++
				}
				rc, err := strconv.Atoi(right[start:pos])
				if err != nil {
					return pd, fmt.Errorf("invalid reroll count: %w", err)
				}
				if rc < 0 {
					return pd, fmt.Errorf("invalid reroll count: %d", rc)
				}
				pd.RerollCount = rc
			}
			continue
		case '>', '<', '=', '!':
			// success operator: two-char first (>=, <=, !=)
			if pos+1 < len(right) {
				two := right[pos : pos+2]
				if two == ">=" || two == "<=" || two == "!=" {
					pos += 2
					if pos >= len(right) || right[pos] < '0' || right[pos] > '9' {
						return pd, fmt.Errorf("parse error: missing success value in %q", s)
					}
					start := pos
					for pos < len(right) && right[pos] >= '0' && right[pos] <= '9' {
						pos++
					}
					sv, err := strconv.Atoi(right[start:pos])
					if err != nil {
						return pd, fmt.Errorf("invalid success value: %w", err)
					}
					pd.SuccessOp = two
					pd.SuccessVal = sv
					continue
				}
			}
			// single-char success op
			if right[pos] == '>' || right[pos] == '<' || right[pos] == '=' {
				op := right[pos]
				pos++
				if pos >= len(right) || right[pos] < '0' || right[pos] > '9' {
					return pd, fmt.Errorf("parse error: missing success value in %q", s)
				}
				start := pos
				for pos < len(right) && right[pos] >= '0' && right[pos] <= '9' {
					pos++
				}
				sv, err := strconv.Atoi(right[start:pos])
				if err != nil {
					return pd, fmt.Errorf("invalid success value: %w", err)
				}
				pd.SuccessOp = string(op)
				pd.SuccessVal = sv
				continue
			}
			return pd, fmt.Errorf("unexpected token %q in %q", string(right[pos]), s)
		case '+', '-', 'x', '*', '/':
			// arithmetic modifier -- consumes the rest and returns
			op := right[pos]
			if op == '+' || op == '-' || op == 'x' || op == '*' || op == '/' {
				pd.ModFunc = string(op)
				if pd.ModFunc == "*" {
					pd.ModFunc = "x"
				}
				pos++
				if pos >= len(right) {
					return pd, fmt.Errorf("parse error: missing modifier value in %q", s)
				}
				mv, err := strconv.Atoi(right[pos:])
				if err != nil {
					return pd, fmt.Errorf("invalid modifier value: %w", err)
				}
				pd.ModVal = mv
				pos = len(right)
				break
			}
		default:
			return pd, fmt.Errorf("unexpected token %q in %q", string(right[pos]), s)
		}
	}

	return pd, nil
}

// RollParsed rolls a ParsedDice using the provided rng. If rng is nil the package
// default RNG is used. Returns a RollResult where AllRolls are per-die totals
// (after exploding) and Rolls are the final contributing values after keep/drop.
func RollParsed(pd ParsedDice, rng *rand.Rand) (RollResult, error) {
	if rng == nil {
		rng = defaultRand
	}
	var res RollResult
	if pd.Count < 0 {
		return res, fmt.Errorf("negative dice count %d", pd.Count)
	}
	if pd.Count > MaxDiceCount {
		return res, fmt.Errorf("dice count %d exceeds limit %d", pd.Count, MaxDiceCount)
	}

	// roll each die, allowing exploding to add to that die's total
	totalRollCalls := 0
	// helper to evaluate reroll predicate on a face value
	matchReroll := func(face int) bool {
		if pd.RerollVal == 0 {
			return false
		}
		switch pd.RerollOp {
		case "=":
			return face == pd.RerollVal
		case "!=":
			return face != pd.RerollVal
		case ">":
			return face > pd.RerollVal
		case "<":
			return face < pd.RerollVal
		case ">=":
			return face >= pd.RerollVal
		case "<=":
			return face <= pd.RerollVal
		default:
			return face == pd.RerollVal
		}
	}
	switch pd.Type {
	case "F":
		// Fate: generate 1..3 then convert to -1..1
		for i := 0; i < pd.Count; i++ {
			// roll
			totalRollCalls++
			if totalRollCalls > MaxTotalRolls {
				return res, fmt.Errorf("exceeded max roll limit")
			}
			die := rng.IntN(pd.Sides) + 1
			// apply reroll predicate if requested (works on the raw face value)
			if matchReroll(die) {
				if pd.RerollOnce {
					// reroll once
					totalRollCalls++
					if totalRollCalls > MaxTotalRolls {
						return res, fmt.Errorf("exceeded max roll limit")
					}
					die = rng.IntN(pd.Sides) + 1
					res.RerollsPerformed++
				} else {
					// reroll until predicate false or per-die limit reached
					per := 0
					for matchReroll(die) {
						if pd.RerollCount > 0 && per >= pd.RerollCount {
							break
						}
						totalRollCalls++
						if totalRollCalls > MaxTotalRolls {
							return res, fmt.Errorf("exceeded max roll limit")
						}
						die = rng.IntN(pd.Sides) + 1
						res.RerollsPerformed++
						per++
					}
				}
			}
			adj := die - 2
			res.AllRolls = append(res.AllRolls, adj)
		}
	case "d":
		if pd.Sides <= 0 {
			return res, fmt.Errorf("invalid sides %d", pd.Sides)
		}
		for i := 0; i < pd.Count; i++ {
			// initial roll
			totalRollCalls++
			if totalRollCalls > MaxTotalRolls {
				return res, fmt.Errorf("exceeded max roll limit")
			}
			die := rng.IntN(pd.Sides) + 1

			// apply reroll predicate before exploding
			if matchReroll(die) {
				if pd.RerollOnce {
					totalRollCalls++
					if totalRollCalls > MaxTotalRolls {
						return res, fmt.Errorf("exceeded max roll limit")
					}
					die = rng.IntN(pd.Sides) + 1
					res.RerollsPerformed++
				} else {
					per := 0
					for matchReroll(die) {
						if pd.RerollCount > 0 && per >= pd.RerollCount {
							break
						}
						totalRollCalls++
						if totalRollCalls > MaxTotalRolls {
							return res, fmt.Errorf("exceeded max roll limit")
						}
						die = rng.IntN(pd.Sides) + 1
						res.RerollsPerformed++
						per++
					}
				}
			}

			dieTotal := die
			if pd.Explode {
				// keep exploding while we hit the maximum face
				for die == pd.Sides {
					totalRollCalls++
					if totalRollCalls > MaxTotalRolls {
						return res, fmt.Errorf("exceeded max roll limit")
					}
					die = rng.IntN(pd.Sides) + 1
					dieTotal += die
				}
			}
			res.AllRolls = append(res.AllRolls, dieTotal)
		}
	default:
		return res, fmt.Errorf("unsupported dice type %q", pd.Type)
	}

	// apply keep/drop
	keep := make([]int, 0, len(res.AllRolls))
	if pd.KeepDropAction == "" {
		keep = append(keep, res.AllRolls...)
	} else {
		// clamp count
		k := pd.KeepDropCount
		if k < 0 {
			k = 0
		}
		if k > len(res.AllRolls) {
			k = len(res.AllRolls)
		}

		// build index-sorted view
		type pair struct{ v, i int }
		ps := make([]pair, 0, len(res.AllRolls))
		for i, v := range res.AllRolls {
			ps = append(ps, pair{v: v, i: i})
		}
		if pd.KeepDropAction == "k" {
			// keep highest or lowest
			if pd.KeepDropWhich == "l" {
				// keep lowest k
				sort.Slice(ps, func(i, j int) bool {
					if ps[i].v == ps[j].v {
						return ps[i].i < ps[j].i
					}
					return ps[i].v < ps[j].v
				})
				keepIdx := map[int]struct{}{}
				for i := 0; i < k; i++ {
					keepIdx[ps[i].i] = struct{}{}
				}
				for i := 0; i < len(res.AllRolls); i++ {
					if _, ok := keepIdx[i]; ok {
						keep = append(keep, res.AllRolls[i])
					}
				}
			} else {
				// keep highest k
				sort.Slice(ps, func(i, j int) bool {
					if ps[i].v == ps[j].v {
						return ps[i].i < ps[j].i
					}
					return ps[i].v > ps[j].v
				})
				keepIdx := map[int]struct{}{}
				for i := 0; i < k; i++ {
					keepIdx[ps[i].i] = struct{}{}
				}
				for i := 0; i < len(res.AllRolls); i++ {
					if _, ok := keepIdx[i]; ok {
						keep = append(keep, res.AllRolls[i])
					}
				}
			}
		} else if pd.KeepDropAction == "d" {
			// drop highest or lowest k
			if pd.KeepDropWhich == "h" {
				// drop highest k -> keep the rest
				sort.Slice(ps, func(i, j int) bool {
					if ps[i].v == ps[j].v {
						return ps[i].i < ps[j].i
					}
					return ps[i].v > ps[j].v
				})
				dropIdx := map[int]struct{}{}
				for i := 0; i < k; i++ {
					dropIdx[ps[i].i] = struct{}{}
				}
				for i := 0; i < len(res.AllRolls); i++ {
					if _, ok := dropIdx[i]; !ok {
						keep = append(keep, res.AllRolls[i])
					}
				}
			} else {
				// drop lowest k
				sort.Slice(ps, func(i, j int) bool {
					if ps[i].v == ps[j].v {
						return ps[i].i < ps[j].i
					}
					return ps[i].v < ps[j].v
				})
				dropIdx := map[int]struct{}{}
				for i := 0; i < k; i++ {
					dropIdx[ps[i].i] = struct{}{}
				}
				for i := 0; i < len(res.AllRolls); i++ {
					if _, ok := dropIdx[i]; !ok {
						keep = append(keep, res.AllRolls[i])
					}
				}
			}
		}
	}

	res.Rolls = keep

	// compute total from kept rolls
	for _, v := range res.Rolls {
		res.Total += v
	}

	// count successes if requested (applies to kept rolls)
	if pd.SuccessOp != "" {
		for _, v := range res.Rolls {
			ok := false
			switch pd.SuccessOp {
			case ">=":
				ok = v >= pd.SuccessVal
			case "<=":
				ok = v <= pd.SuccessVal
			case ">":
				ok = v > pd.SuccessVal
			case "<":
				ok = v < pd.SuccessVal
			case "=":
				ok = v == pd.SuccessVal
			}
			if ok {
				res.Successes++
			}
		}
	}

	// record total roll calls
	res.TotalRollCalls = totalRollCalls

	// Apply modifier to the total
	switch pd.ModFunc {
	case "":
		// nothing
	case "+":
		res.Total += pd.ModVal
	case "-":
		res.Total -= pd.ModVal
	case "x":
		res.Total *= pd.ModVal
	case "/":
		if pd.ModVal == 0 {
			return res, fmt.Errorf("division by zero")
		}
		res.Total /= pd.ModVal
	default:
		return res, fmt.Errorf("unsupported modifier %q", pd.ModFunc)
	}

	return res, nil
}
