package feature

import (
	"math/big"
	"crypto/rand"
	log "github.com/sirupsen/logrus"
	"strings"
	"strconv"
	"YmirBot/cmd/db"
	"encoding/json"
	"time"
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

func RollDiceModifierWithHistory(numDice int64, numSides int64, modifier int64, database db.Database, name string) *big.Int {
	var dice []Dice

	for i := int64(0); i < numDice; i++ {
		dice = append(dice, Dice{big.NewInt(numSides)})
	}

	total := big.NewInt(0)

	for _, die := range dice {
		roll_result := die.Roll()
		total.Add(total, roll_result)

		if numSides == int64(20) {
			key := name+"/20"
			data, err := database.Get(key)
			if err != nil {
				roll_results := make([]RollResult, 1)
				roll_results[0] = RollResult{time.Now(), roll_result.Int64()}
				new_record := PersonHistory{name, roll_results}

				value, err := json.Marshal(new_record)
				if err != nil {
					log.Error("Could not marshall data to DB: ", err)
				}

				database.Put(key, value)
			} else {
				var previous_record PersonHistory
				err := json.Unmarshal(data, previous_record)
				if err != nil {
					log.Error("Could not unmarshal data from DB: ", err)
				}

				previous_record.RollResults = append(previous_record.RollResults, RollResult{time.Now(), roll_result.Int64()})
				value, err := json.Marshal(previous_record)
				if err != nil {
					log.Error("Could not marshall data to DB: ", err)
				}

				err = database.Put(key, value)
				if err != nil {
					log.Error("Issue with DB: ", err)
				}
			}
		}
	}

	return total.Add(total, big.NewInt(modifier))
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