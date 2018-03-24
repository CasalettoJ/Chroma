package blockchain

import (
	"os"
)

// DoesDBExist checks for an existing blockchain db file
func DoesDBExist() bool {
	_, err := os.Stat(DBdbfile)
	if os.IsNotExist(err) {
		return false
	}

	CheckAnxiety(err)
	return true
}
