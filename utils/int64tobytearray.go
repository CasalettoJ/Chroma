package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

// Int64ToByteArray returns []byte big endian order representation of an in64
func Int64ToByteArray(num int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
