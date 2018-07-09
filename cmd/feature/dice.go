package feature

import (
	"math/big"
	"crypto/rand"
	log "github.com/sirupsen/logrus"
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

func RollDice(numDice int64, numSides int64) *big.Int {

	var dice []Dice

	for i := int64(0); i < numDice; i++ {
		dice = append(dice, Dice{big.NewInt(numSides)})
	}

	total := big.NewInt(0)

	for _, die := range dice {
		total.Add(total, die.Roll())
	}

	return total
}

func RollDiceModifier(numDice int64, numSides int64, modifier int64) *big.Int {
	original := RollDice(numDice, numSides)

	return original.Add(original, big.NewInt(modifier))
}