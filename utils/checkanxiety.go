package utils

import "log"

// CheckAnxiety detects if we gotta panic.
func CheckAnxiety(err error) {
	if err != nil {
		log.Panic(err)
	}
}
