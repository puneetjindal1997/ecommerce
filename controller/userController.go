package controller

import (
	"ecommerce/constant"
	"ecommerce/database"
	"ecommerce/helper"
	"ecommerce/types"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

/*
 *	Email verify process
 *
 */
func VerifyEmail(c *gin.Context) {
	var req types.Verification
	postBodyErr := c.BindJSON(&req)
	if postBodyErr != nil {
		log.Println(postBodyErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": postBodyErr.Error()})
		return
	}
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EmailValidationError})
		return
	}
	resp := database.Mgr.GetSingleRecordByEmail(req.Email, constant.VerificationsCollection)
	if resp.Otp != 0 {
		sec := resp.CreatedAt + constant.OtpValidation
		// checking if otp expire
		if sec < time.Now().Unix() {
			req, checkEmail := helper.SendEmailSendGrid(req)
			if checkEmail != nil {
				log.Println(postBodyErr)
				c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EmailValidationError})
				return
			}
			// update email send time
			req.CreatedAt = time.Now().Unix()
			// update operation
			database.Mgr.UpdateVerification(req, constant.VerificationsCollection)
			c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OptAlreadySentError})
			return
		}
	}
	req, checkEmail := helper.SendEmailSendGrid(req)
	if checkEmail != nil {
		log.Println(postBodyErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EmailValidationError})
		return
	}
	req.CreatedAt = time.Now().Unix() // unix timestamp

	// insertion of record
	database.Mgr.Insert(req, constant.VerificationsCollection)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
}

/*
 *	Verify otp end to end
 *
 */
func VerifyOtp(c *gin.Context) {
	var req types.Verification
	postBodyErr := c.BindJSON(&req)
	if postBodyErr != nil {
		log.Println(postBodyErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": postBodyErr.Error()})
		return
	}

	// checking email
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EmailValidationError})
		return
	}
	fmt.Println(req.Email, req.Otp)
	// checking otp field not be empty
	if req.Otp <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OtpValidationError})
		return
	}

	// checking the records in verification collection
	resp := database.Mgr.GetSingleRecordByEmail(req.Email, constant.VerificationsCollection)
	// if status or email is already verified
	if resp.Status {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.AlreadyVerifiedError})
		return
	}
	sec := resp.CreatedAt + constant.OtpValidation
	// otp coming in the request and from db can't matched
	if resp.Otp != req.Otp {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OtpValidationError})
		return
	}
	// otp expired
	if sec < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OtpExpiredValidationError})
		return
	}
	// if all good then we will verified the email
	req.Status = true
	req.CreatedAt = time.Now().Unix()
	err := database.Mgr.UpdateEmailVerifiedStatus(req, constant.VerificationsCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OtpValidationError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
}

/*
 *	Controller related to user registeration
 *
 */
func RegisterUser(c *gin.Context) {
	var userClient types.UserClient
	var dbUser types.User

	// binding the payload
	reqErr := c.BindJSON(&userClient)
	if reqErr != nil {
		log.Println(reqErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": reqErr.Error()})
		return
	}
	// payload error handle
	err := helper.CheckUserValidation(userClient)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	// check if email is verified
	resp := database.Mgr.GetSingleRecordByEmail(userClient.Email, constant.VerificationsCollection)
	if !resp.Status {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EmailIsNotVerified})
		return
	}

	// check for duplicate user
	respUser := database.Mgr.GetSingleRecordByEmailForUser(userClient.Email, constant.UserCollection)
	if respUser.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.AlreadyRegisterWithThisEmail})
		return
	}
	dbUser.Email = userClient.Email
	dbUser.Name = userClient.Name
	dbUser.Phone = userClient.Phone
	// encrypted password genration
	encryptedPass := helper.GenPassHash(userClient.Password)
	dbUser.Password = encryptedPass
	dbUser.CreatedAt = time.Now().Unix()
	dbUser.UpdatedAt = time.Now().Unix()
	err = database.Mgr.Insert(dbUser, constant.UserCollection)
	if err != nil {
		log.Println(reqErr)
		// custom error return todo
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err})
		return
	}
	dbUser.Password = ""
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": dbUser})
}
