package dnd

import (
	"HelloBot/cmd/db"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ----------------------------------------------------------------------------
// Dice statistics history.
// ----------------------------------------------------------------------------

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

// ----------------------------------------------------------------------------
// Group roll history.
// ----------------------------------------------------------------------------

type GroupRoll struct {
	GroupRollMap map[string][]string `json:"group_roll_map"`
}

func GroupStart(database db.Database) string {
	value, err := json.Marshal(GroupRoll{make(map[string][]string)})
	
	if err != nil {
		log.Error("Could not marshall data to DB: ", err)
	}
	
	database.Put("group_roll", value)
	return "Beginning recording of a group of rolls."
}
	
func GroupSetRoll(name string, rollNum *big.Int, database db.Database) string {
	result := "recorded"

	data, err := database.Get("group_roll")
	if err != nil {
		return "Not currently recording group rolls."
	}
	
	var groupRollStruct GroupRoll
	err = json.Unmarshal(data, &groupRollStruct)
	if err != nil {
		log.Error("Could not unmarshal data from DB: ", err)
	}
	groupRollMap := groupRollStruct.GroupRollMap
	
	bucket := ""
	for k, v := range groupRollMap {
		for i, stored := range v {
			if (stored == name) {
				bucket = k

				v[i] = v[len(v) - 1]
				groupRollMap[k] = v[:len(v) - 1]
			
				result = "overriden"
			}
		}	
    }
	if (bucket != "" && len(groupRollMap[bucket]) == 0) {
		delete(groupRollMap, bucket)
	}
	
	stringRollNum := strconv.FormatInt(rollNum.Int64(), 10)
	if (groupRollMap[stringRollNum] == nil) {
		groupRollMap[stringRollNum] = []string{name}
	} else {
		groupRollMap[stringRollNum] = append(groupRollMap[stringRollNum], name)
	}
	
	value, err := json.Marshal(groupRollStruct)
	if err != nil {
		log.Error("Could not marshall data to DB: ", err)
	}

	err = database.Put("group_roll", value)
	if err != nil {
		log.Error("Issue with DB: ", err)
	}
	
	return result
}

func GroupEnd(database db.Database) string {
	data, err := database.Get("group_roll")

	if err != nil {
		return "Not currently recording group rolls."
	}

	var currentGroupRoll GroupRoll
	err = json.Unmarshal(data, &currentGroupRoll)
	
	if err != nil {
		return "Not currently recording group rolls."
	}
	
	var rolledNums []int
    for k := range currentGroupRoll.GroupRollMap {
		val, err := strconv.Atoi(k)
		
		if (err == nil) {
			rolledNums = append(rolledNums, val)
		} else {
			log.Error("Issue with converting string to number when ending group: ", err)
		}		
    }
	sort.Sort(sort.Reverse(sort.IntSlice(rolledNums)))

	response := "Group roll results:\n"
	for _, k := range rolledNums {
		stringVal := strconv.Itoa(k)
		response += "\t" + stringVal + ":\t\t\t" +
			strings.Join(currentGroupRoll.GroupRollMap[stringVal], ", ") + "\n"
	}
	
	database.Put("group_roll", nil)
	return response
}