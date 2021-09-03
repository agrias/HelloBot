package nicehash

import (
	"fmt"
	"os/exec"
	"strings"
	"github.com/thedevsaddam/gojsonq/v2"
	"strconv"
)

func GetBalance() float64 {
	argstr := []string{"/Development/workspace/nicehash/nicehash.sh"}
	out, err := exec.Command("/bin/bash", argstr...).Output()

	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed")
	}
	
	jsonstring := strings.ReplaceAll(string(out), "'", "\"")
	jsonstring = strings.ReplaceAll(jsonstring, "True", "true")
	jsonstring = strings.ReplaceAll(jsonstring, "False", "false")

	balance := gojsonq.New().FromString(jsonstring).Find("total.totalBalance")
	numeric_balance, err := strconv.ParseFloat(balance.(string), 64)

	if err == nil {
		fmt.Println(err)
	}

	return numeric_balance
}

func GetPrice() float64 {
	argstr := []string{"https://api.coinbase.com/v2/prices/spot?currency=USD"}
	out, err := exec.Command("/usr/bin/curl", argstr...).Output()

	if err != nil {
		fmt.Println(err)
	}

	jsonstring := string(out)

	price := gojsonq.New().FromString(jsonstring).Find("data.amount")

	numeric_balance, err := strconv.ParseFloat(price.(string), 64)

	if err == nil {
		fmt.Println(err)
	}

	return numeric_balance
}
