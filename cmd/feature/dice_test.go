package feature

import (
	"testing"
	"math/big"
	"fmt"
	"strconv"
	"github.com/stretchr/testify/assert"
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

func t3(one, two, three interface{}) []interface{} {
	return []interface{}{one, two, three}
}

func TestParseDiceString(t *testing.T) {
	string1 := "!roll 1d4"
	string2 := "!roll 1d4+4"
	assert.Equal(t, t3(ParseDiceString(string1)), t3(int64(1), int64(4), int64(0)))
	assert.Equal(t, t3(ParseDiceString(string2)), t3(int64(1), int64(4), int64(4)))
}

func TestFormatRollResults(t *testing.T) {
	array := []*big.Int{big.NewInt(20), big.NewInt(2), big.NewInt(3)}
	name := "test"
	request := "!roll 1d20"

	assert.Equal(t, FormatRollResults(array, name, request), "<@test> \n:star2:   :star2:   **NATURAL 20**   :star2:   :star2:\n20 + *2* = **3**")

	request = "!roll 1d20+1"
	assert.Equal(t, FormatRollResults(array, name, request), "<@test> \n:star2:   :star2:   **NATURAL 20**   :star2:   :star2:\n20 + *2* = **3**")
}