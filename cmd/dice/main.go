package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ppreeper/dice"
)

func main() {
	var dicelist string
	flag.StringVar(&dicelist, "d", "d6", "dice notation")
	flag.Parse()

	dd, err := dice.Parse(dicelist)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	res, err := dice.RollParsed(dd, nil) // use default RNG
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	faces := []string{}
	for _, v := range res.AllRolls {
		faces = append(faces, fmt.Sprintf("%d", v))
	}
	fmt.Println(strings.Join(faces, ";"))
}
