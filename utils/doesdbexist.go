package utils

import (
	"os"

	conf "github.com/casalettoj/chroma/constants"
)

// DoesDBExist checks for an existing blockchain db file
func DoesDBExist() bool {
	_, err := os.Stat(conf.DBdbfile)
	if os.IsNotExist(err) {
		return false
	}

	CheckAnxiety(err)
	return true
}
