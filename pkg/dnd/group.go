package dnd

import (
	"HelloBot/cmd/db"
	"HelloBot/proto"
	"math/big"
	"strings"
)

func ProcessGroupCommand(req *proto.BotRequest, database db.Database) string {
	response := ""
	splits := strings.Split(req.Text, " ")
	
	if (len(splits) == 1) {
		return response;
	}
	
	switch splits[1] {
	case "start":
		response = GroupStart(database)
		
	case "end":
		response = GroupEnd(database)
	
	case "set":
		if (len(splits) > 2) {
			rollAmount := new(big.Int)
			rollAmount, err := rollAmount.SetString(splits[2], 10)
			
			if !err {
				return "Error parsing number for set command."
			}

			var name string

			if len(splits) > 3 {
				name = splits[3]
			} else {
				name = "<@" + req.Name + ">"
			}
			result := GroupSetRoll(name, rollAmount, database)
			response = name + " -- Roll " + result
		}
	
	case "roll":
		if (len(splits) > 2) {
			num_dice, sides, modifier := ParseDiceString(splits[1] + " " + splits[2])
			rolls := RollDiceModifier(num_dice, sides, modifier)

			name := "<@" + req.Name + ">"
			result := GroupSetRoll(name, rolls[len(rolls) - 1], database)
			response = name + " -- Roll " + result
		}
	
	case "help":
		response = "Valid commands...:\n" +
			"\t!group start\n" +
			"\t!group set (number)\n" +
			"\t!group roll (dice)\n" +
			"\t!group end"

	default:
		response = "Invalid group command. \"!group help\" for more details."
	}

	return response;
}