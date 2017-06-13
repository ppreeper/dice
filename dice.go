package main

import (
	"fmt"
	"math/rand"

	"git.preeper.org/peterp/cdn"
)

func main() {
	fmt.Println("Dice")
	fmt.Println(cdn.DieSides("1d6"))
	fmt.Println(Roll(cdn.DieSides("1d6")))
}

//Roll soll the die given
func Roll(sides int) int {
	// rand.Seed(0)
	return rand.Intn(sides)
}
