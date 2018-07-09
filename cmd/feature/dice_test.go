package feature

import (
	"testing"
	"math/big"
	"fmt"
	"strconv"
)

func TestDice_Roll(t *testing.T) {
	d6 := &Dice{big.NewInt(20)}

	results := make(map[string]int)

	for i:= 0; i < 10000; i++ {
		res := d6.Roll()
		fixed := res.String()

		results[fixed] = results[fixed] + 1
	}

	for key, value := range results {
		fmt.Println("Key: "+key+" Value: "+strconv.Itoa(value))
	}
}

func TestRollDice(t *testing.T) {

	for i:= 0; i < 100; i++ {
		fmt.Println(RollDice(100, 20))
	}
}