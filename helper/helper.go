package helper

import (
	"ecommerce/types"
	"errors"
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func ConverStringIntoInt(str string) (int, error) {
	integer, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return integer, err
}

func CheckUserValidation(u types.UserClient) error {
	if u.Email == "" {
		return errors.New("email can't be empty")
	}
	if u.Name == "" {
		return errors.New("name can't be empty")
	}
	if u.Phone == "" {
		return errors.New("phone can't be empty")
	}
	if u.Password == "" {
		return errors.New("password can't be empty")
	}
	return nil
}

func GenPassHash(s string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(bytes)
}
