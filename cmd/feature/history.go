package feature

import (
	"time"
	"math/big"
	"YmirBot/cmd/db"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type PersonHistory struct {
	Name string	`json:"name"`
	Timestamp time.Time `json:"timestamp"`
	RollResults RollResult	`json:"roll_results"`
}

type RollResult struct {
	//Timestamp time.Time	`json:"timestamp"`
	Twenty map[string]int	`json:"twenty"`
	Total int `json:"total"`
}

func makeNewPersonHistory(name string) PersonHistory {

	twenty := map[string]int{
		"0" : 0,
		"1" : 0,
		"2" : 0,
		"3" : 0,
		"4" : 0,
		"5" : 0,
		"6" : 0,
		"7" : 0,
		"8" : 0,
		"9" : 0,
		"10" : 0,
		"11" : 0,
		"12" : 0,
		"13" : 0,
		"14" : 0,
		"15" : 0,
		"16" : 0,
		"17" : 0,
		"18" : 0,
		"19" : 0,
		"20" : 0,
	}

	results := RollResult{twenty, 0}
	return PersonHistory{name, time.Now(), results}
}

func increment20map(twenty map[string]int, roll int64) {
	twenty[strconv.FormatInt(roll, 10)] += 1
}

func ModifyDiceHistory(roll_result []*big.Int, database db.Database, name string) {
	log.Infof("Modifying %s\n", name)

	key := name+"/20"
	data, err := database.Get(key)
	if err != nil {
		new_history := makeNewPersonHistory(name)

		for i := 0; i < len(roll_result); i++ {
			increment20map(new_history.RollResults.Twenty, roll_result[i].Int64())
			new_history.RollResults.Total += 1
		}

		value, err := json.Marshal(new_history)
		if err != nil {
			log.Error("Could not marshall data to DB: ", err)
		}

		database.Put(key, value)
	} else {
		var previous_record PersonHistory
		err := json.Unmarshal(data, &previous_record)
		if err != nil {
			log.Error("Could not unmarshal data from DB: ", err)
		}

		for i := 0; i < len(roll_result); i++ {
			increment20map(previous_record.RollResults.Twenty, roll_result[i].Int64())
			previous_record.RollResults.Total += 1
		}

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

func GetDiceHistoryStats(database db.Database, name string) (string) {
	key := name+"/20"
	data, err := database.Get(key)
	if err != nil {
		return "No data"
	} else {
		var previous_record PersonHistory
		err := json.Unmarshal(data, &previous_record)
		if err != nil {
			log.Error("Could not unmarshal data from DB: ", err)
		}

		json_text, err := json.MarshalIndent(previous_record.RollResults.Twenty, "", "")
		if err != nil {
			panic(err)
		}

		return string(json_text)
	}
}
