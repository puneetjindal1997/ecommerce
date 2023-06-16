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

func CheckProductValidation(p types.ProductClient) error {
	if p.Description == "" {
		return errors.New("description of product can't be empty")
	}
	if p.Name == "" {
		return errors.New("name of product can't be empty")
	}
	if p.ImageUrl == "" {
		return errors.New("image of product can't be empty")
	}
	if p.Price == 0 {
		return errors.New("price of product can't be 0")
	}
	return nil
}
