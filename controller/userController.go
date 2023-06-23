package controller

import (
	"ecommerce/auth"
	"ecommerce/constant"
	"ecommerce/database"
	"ecommerce/helper"
	"ecommerce/types"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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
	dbUser.UserType = constant.NormalUser
	// encrypted password genration
	encryptedPass := helper.GenPassHash(userClient.Password)
	dbUser.Password = encryptedPass
	dbUser.CreatedAt = time.Now().Unix()
	dbUser.UpdatedAt = time.Now().Unix()
	id, err := database.Mgr.Insert(dbUser, constant.UserCollection)
	if err != nil {
		log.Println(reqErr)
		// custom error return todo
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err})
		return
	}

	// jwt struct prepare
	jwtWrapper := auth.JwtWrapper{
		ScretKey:       os.Getenv("JwtSecrets"),
		Issuer:         os.Getenv("JwtIssuer"),
		ExpirationTime: 48,
	}
	userId := id.(primitive.ObjectID)
	// gen token
	token, err := jwtWrapper.GenrateToken(userId, userClient.Email, dbUser.UserType)
	if err != nil {
		log.Println(err)
		// custom error return todo
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	dbUser.Password = ""
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": dbUser, "token": token})
}

func UserLogin(c *gin.Context) {
	var loginReq types.Login
	err := c.BindJSON(&loginReq)
	if err != nil {
		log.Println(err)
		// custom error return todo
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	user := database.Mgr.GetSingleRecordByEmailForUser(loginReq.Email, constant.UserCollection)
	if user.Email == "" {
		log.Println(err)
		// custom error return todo
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.NotRegisteredUser})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		log.Println(err)
		// custom error return todo
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.PasswordNotMatchedError})
		return
	}

	jwtWrapper := auth.JwtWrapper{
		ScretKey:       os.Getenv("JwtSecrets"),
		Issuer:         os.Getenv("JwtIssuer"),
		ExpirationTime: 48,
	}

	token, err := jwtWrapper.GenrateToken(user.Id, user.Email, user.UserType)
	if err != nil {
		log.Println(err)
		// custom error return todo
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func AddToCart(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EmailIsNotVerified})
		return
	}
	userDbResp := database.Mgr.GetSingleRecordByEmail(email.(string), constant.UserCollection)
	if userDbResp.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.UserDoesNotExists})
		return
	}
	address, err := database.Mgr.GetSingleAddress(userDbResp.ID, constant.AddressCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	if address.Address1 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.AddressNotExists})
		return
	}
	var cart types.CartClient
	var cartDb types.Cart
	err = c.BindJSON(&cart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	projectId, _ := primitive.ObjectIDFromHex(cart.ProductID)
	userId, _ := primitive.ObjectIDFromHex(cart.UserID)
	cartDb.ProductID = projectId
	cartDb.UserID = userId
	cartDb.Checkout = false
	_, err = database.Mgr.Insert(cartDb, constant.CartCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "error": false})
}

func AddAddressOfUser(c *gin.Context) {
	var addressReq types.AddressClient
	err := c.BindJSON(&addressReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	userId, err := primitive.ObjectIDFromHex(addressReq.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	var addressDb types.Address

	addressDb.Address1 = addressReq.Address1
	addressDb.UserId = userId
	addressDb.City = addressReq.City
	addressDb.Country = addressReq.Country

	_, err = database.Mgr.Insert(addressDb, constant.AddressCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "error": false})
}

func GetSingleUser(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	user := database.Mgr.GetSingleUserByUserId(userId, constant.UserCollection)
	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.UserDoesNotExists})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "success", "error": false, "data": user})
}

func UpdateUser(c *gin.Context) {
	var userUpdate types.UserUpdateClient
	err := c.BindJSON(&userUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	userId, err := primitive.ObjectIDFromHex(userUpdate.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	userResp := database.Mgr.GetSingleUserByUserId(userId, constant.UserCollection)
	if userResp.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.UserDoesNotExists})
		return
	}
	var user types.User
	user.Id = userId
	user.Email = userResp.Email
	user.Password = userResp.Password
	user.Phone = userResp.Phone
	user.UserType = userResp.UserType
	user.UpdatedAt = time.Now().Unix()
	user.CreatedAt = userResp.CreatedAt

	if userUpdate.Email != "" {
		user.Email = userUpdate.Email
	}
	if userUpdate.Password != "" {
		user.Password = helper.GenPassHash(userUpdate.Password)
	}
	if userUpdate.Phone != "" {
		user.Phone = userUpdate.Phone
	}
	if userUpdate.Name != "" {
		user.Name = userUpdate.Name
	}

	err = database.Mgr.UpdateUser(user, constant.UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "error": false})
}
