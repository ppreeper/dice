# dice

![Build Status](https://github.com/ppreeper/dice/actions/workflows/go.yml/badge.svg)
![Coverage Status](badges/coverage.svg)

A dice rolling library for use in games. Specifically to be used for games using [dice notation](https://en.wikipedia.org/wiki/Dice_notation).

Usage
-----

This package provides a small, safe parser and pure roll function. The API is intentionally split:

- Parse(notation string) (ParsedDice, error) — parses dice notation into a value.
- RollParsed(pd ParsedDice, rng *rand.Rand) (RollResult, error) — rolls the parsed dice using a provided RNG. Pass nil to use the package default RNG.

Supported notation (subset):

- NdS, dS (e.g. `3d6`, `d20`)
- Fate dice: `NF` (e.g. `4F` produces values in -1..1 per die)
- Modifiers: `+ - x /` (x is multiplication; `*` is accepted and normalized to `x`)
- Exploding: `!` after sides, e.g. `2d6!`
- Keep/Drop: `kh`/`kl`/`dh`/`dl` or `k`/`d` shorthand, with a count. Examples: `4d6kh3`, `4d6d1`, `4d6k3`.
- Percentile: `d%` is accepted as `d100`.
 - Rerolls: `r` / `ro` with comparators and per-die cap. Examples:
   - `r1` (reroll faces equal to 1 until they are different)
   - `ro1` (reroll once if equal to 1)
   - comparator forms: `r<2`, `r>=5`, `r!=3`
   - per-die cap: append `#N` to limit rerolls per die, e.g. `r1#2`
 - Success counting: `>=`, `>`, `<=`, `<`, `=` after the dice term, e.g. `10d10>=8`.

Limits & safety
---------------

To avoid resource exhaustion or accidental infinite loops the parser/runner enforces reasonable limits:

- Max dice count: 1000
- Max sides per die: 1,000,000
- Max total RNG calls per roll (including explosions): 10,000

Exploding on a 1-sided die is rejected to prevent infinite exploding loops.

Examples
--------

Basic roll:

```go
pd, err := dice.Parse("3d6")
if err != nil { /* handle */ }
res, err := dice.RollParsed(pd, nil) // use default RNG
// res.Rolls -> per-die results
// res.Total -> sum
```

Deterministic tests (seed RNG):

```go
// rand/v2 PCG-based RNG for deterministic results
rng := rand.New(rand.NewPCG(600, 601))
pd, _ := dice.Parse("4d6kh3")
res, _ := dice.RollParsed(pd, rng)
// res contains deterministic results
```

Exploding dice and keep-highest example:

```go
pd, _ := dice.Parse("4d6!kh3")
res, _ := dice.RollParsed(pd, nil)
```

Rerolls and success counting
----------------------------

```go
// reroll ones until not 1
pd, _ := dice.Parse("4d6r1")
res, _ := dice.RollParsed(pd, nil)

// reroll once (ro)
pd2, _ := dice.Parse("4d6ro1")
res2, _ := dice.RollParsed(pd2, nil)

// success counting: count how many d10 results are >= 8
pd3, _ := dice.Parse("10d10>=8")
res3, _ := dice.RollParsed(pd3, nil)
// res3.Successes contains the number of successes
```

Grammar & token ordering
------------------------

After the `dS` portion (and optional `!` for exploding) tokens may appear in any order. Supported trailing tokens are:

- Keep/Drop: `k` or `d` optionally followed by `h`/`l` and a count, e.g. `kh3` or `d2`.
- Reroll: `r` or `ro` plus comparator and value. Comparator can be `=`, `!=`, `<`, `<=`, `>`, `>=`. Example: `r<2`, `ro!=3`, `r1` (shorthand for `r=1`). Optionally append per-die cap `#N` (e.g. `r1#2`).
- Success operator: `>=N`, `>N`, `<=N`, `<N`, `=N` (counts successes among kept dice).
- Arithmetic modifier: `+N`, `-N`, `xN`.

Examples of equivalent orderings:

- `4d6kh3r1` == `4d6r1kh3`
- `2d6!r<2+1` == `2d6!+1r<2`

Notes
-----
- Rerolls are applied to the raw die face before explosions. Exploding is evaluated on the final face after rerolls.
- All RNG calls (initial, rerolls, explosions) are counted toward MaxTotalRolls. If that limit is exceeded RollParsed returns an error.
- Parse errors now include the original notation in messages for easier debugging.

Contributing & extending
-------------------------

The parser is intentionally conservative and we welcome contributions — please open an issue or PR. Implemented extensions include reroll comparator predicates (r/ro with `=`, `!=`, `<`, `<=`, `>`, `>=`), per-die reroll caps using `#N` (e.g. `r1#2`), and success thresholds (`>=`, `>`, `<=`, `<`, `=`).

If you want additional features (compound/penetrating explosions, multi-term expression parsing, advanced reroll predicates, a CLI/HTTP API, or analytic helpers), open an issue or submit a PR — we can extend the parser and the pure roller incrementally.
