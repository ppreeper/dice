package main

import (
	"fmt"

	_ "dice/cdn"
)

func main() {
	fmt.Println("Dice")
	fmt.Println(roll.Roll(6))
}
