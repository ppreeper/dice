package main

import (
	"fmt"
	"math/rand"

	"git.preeper.org/peterp/cdn"
)

func main() {
	fmt.Println("Dice")
	fmt.Println(cdn.DieCount("1d6"))
	fmt.Println(cdn.DieSides("1d6"))
	fmt.Println(Roll(6))
	fmt.Println(Roll(6))
	fmt.Println(Roll(6))
	fmt.Println(Roll(6))
	fmt.Println(Roll(6))
	fmt.Println(Roll(6))
}

//Roll soll the die given
func Roll(sides int) int {
	// rand.Seed(0)
	return rand.Intn(sides)
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
