package controller

import (
	"ecommerce/constant"
	"ecommerce/database"
	"ecommerce/helper"
	"ecommerce/types"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterProduct(c *gin.Context) {
	userEmail, ok := c.Get("email")
	fmt.Println(userEmail, ok)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": constant.NotRegisteredUser})
		return
	}
	user := database.Mgr.GetSingleRecordByEmailForUser(userEmail.(string), constant.UserCollection)
	if user.UserType != constant.AdminUser {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": constant.NotAuthorizedUserError})
		return
	}
	var productRequest types.ProductClient
	var p types.Product
	err := c.BindJSON(&productRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}
	err = helper.CheckProductValidation(productRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}
	p.Name = productRequest.Name
	p.Description = productRequest.Description
	p.ImageUrl = productRequest.ImageUrl
	p.Price = productRequest.Price
	p.MetaInfo = productRequest.MetaInfo
	p.CreatedAt = time.Now().Unix()
	p.UpdatedAt = time.Now().Unix()
	id, err := database.Mgr.Insert(p, constant.ProductCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}
	p.ID = id.(primitive.ObjectID)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": p})
}
