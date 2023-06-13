package helper

import (
	"log"
	"strconv"
)

func ConverStringIntoInt(str string) (int, error) {
	integer, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return integer, err
}
