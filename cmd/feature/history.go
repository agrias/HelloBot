package feature

import "time"

type PersonHistory struct {
	Name string	`json:"name"`
	RollResults []RollResult	`json:"roll_results"`
}

type RollResult struct {
	Timestamp time.Time	`json:"timestamp"`
	Twenty int64	`json:"twenty"`
}