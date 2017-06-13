package roll

import "math/rand"

//Roll soll the die given
func Roll(sides int) int {
	// rand.Seed(0)
	return rand.Intn(sides)
}
