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
 *	resend otp if email exists
 *
 */
func ResendOTPEmail(c *gin.Context) {
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
	if resp.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EmailValidationError})
		return
	}
	req, checkEmail := helper.SendEmailSendGrid(req)
	if checkEmail != nil {
		log.Println(postBodyErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EmailValidationError})
		return
	}
	req.CreatedAt = time.Now().Unix()
	database.Mgr.UpdateVerification(req, constant.VerificationsCollection)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
}
