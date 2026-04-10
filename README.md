# dice

![Build Status](badges/ci.svg)
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
rng := rand.New(rand.NewSource(600))
pd, _ := dice.Parse("4d6kh3")
res, _ := dice.RollParsed(pd, rng)
// res contains deterministic results
```

Exploding dice and keep-highest example:

```go
pd, _ := dice.Parse("4d6!kh3")
res, _ := dice.RollParsed(pd, nil)
```

Contributing & extending
-------------------------

The parser is intentionally conservative. If you want more features (rerolls, success thresholds, compound exploding, etc.) open an issue or PR — we can extend the parser and the pure roller incrementally.
