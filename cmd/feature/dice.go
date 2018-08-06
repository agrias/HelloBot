package feature

import (
	"math/big"
	"crypto/rand"
	log "github.com/sirupsen/logrus"
	"strings"
	"strconv"
	"YmirBot/cmd/db"
)

type Dice struct {
	Sides *big.Int
}

func (d *Dice) Roll() *big.Int {
	value, err := rand.Int(rand.Reader, d.Sides)
	if err != nil {
		log.Error("Something wrong with ")
	}

	return value.Add(value, big.NewInt(1))
}

// Returns an array containing the results of each individual roll, and the
// last element containing the total of the rolls.
func RollDice(numDice int64, numSides int64) []*big.Int {
	var dice []Dice

	for i := int64(0); i < numDice; i++ {
		dice = append(dice, Dice{big.NewInt(numSides)})
	}

	var allRolls []*big.Int
	total := big.NewInt(0)

	for _, die := range dice {
		rollAmount := die.Roll()
		allRolls = append(allRolls, die.Roll())
		total.Add(total, rollAmount)
	}
	allRolls = append(allRolls, total)

	return allRolls
}

// Returns an array containing the results of each individual roll, and the
// last two elements containing the modifier and total of the rolls respectively.
func RollDiceModifier(numDice int64, numSides int64, modifier int64) []*big.Int {
	allRolls := RollDice(numDice, numSides)

	lastIndex := len(allRolls) - 1
	total := allRolls[lastIndex]

	allRolls[lastIndex] = big.NewInt(modifier)
	total.Add(total, allRolls[lastIndex])
	allRolls = append(allRolls, total)

	return allRolls
}

// Returns an array containing the results of each individual roll, and the
// last two elements containing the modifier and total of the rolls respectively.
func RollDiceModifierWithHistory(numDice int64, numSides int64, modifier int64, database db.Database, name string) []*big.Int {
	var dice []Dice

	for i := int64(0); i < numDice; i++ {
		dice = append(dice, Dice{big.NewInt(numSides)})
	}

	total := big.NewInt(0)

	var roll_values []*big.Int

	for _, die := range dice {
		roll_result := die.Roll()
		total.Add(total, roll_result)
		roll_values = append(roll_values, roll_result)
	}

	if numSides == int64(20) {
		ModifyDiceHistory(roll_values, database, name)
	}
	
	bigIntModifier := big.NewInt(modifier)
	total.Add(total, bigIntModifier)
	roll_values = append(roll_values, bigIntModifier)
	roll_values = append(roll_values, total)

	return roll_values
}

func ParseDiceString(text string) (int64, int64, int64) {
	splits := strings.Split(text, " ")
	numbers := strings.Split(splits[1], "d")
	modifier := strings.Split(numbers[1], "+")

	if len(modifier) != 2 {
		return StringToInt64(numbers[0]), StringToInt64(numbers[1]), 0
	}

	return StringToInt64(numbers[0]), StringToInt64(modifier[0]), StringToInt64(modifier[1])
}

func StringToInt64(text string) (int64) {
	value, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		panic(err)
	}

	return value
}

func FormatRollResults(rollResults []*big.Int, name string, request string) string {
	displayString := "<@" + name + "> "
	numResults := len(rollResults)
	
	// Extra formatting for natural 1s and 20s.
	splits := strings.Split(request, " ")
	if (splits[1] == "1d20") {
		if (rollResults[0].Cmp(big.NewInt(20)) == 0) {
			displayString = displayString + "\n:star2:   :star2:   **NATURAL 20**   :star2:   :star2:\n"
		} else if (rollResults[0].Cmp(big.NewInt(1)) == 0) {
			displayString = displayString + "\n:zap:   :zap:   **NATURAL 1**   :zap:   :zap:\n"
		}
	}

	// Adding all the actual rolls from the reuslts.
	for _, result := range rollResults[:numResults - 2] {
		displayString = displayString + result.String() + " + "
	}
	
	// Adding the modifier and total.
	displayString = displayString + "*" + rollResults[numResults - 2].String() + "* = "
	displayString = displayString + "**" + rollResults[numResults - 1].String() + "**"

	return displayString
}