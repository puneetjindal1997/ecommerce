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

func ListProductsController(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	pageInt, _ := helper.ConverStringIntoInt(page)
	limitInt, _ := helper.ConverStringIntoInt(limit)
	offsetInt, _ := helper.ConverStringIntoInt(offset)

	dbResp, count, err := database.Mgr.GetListProducts(pageInt, limitInt, offsetInt, constant.ProductCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": map[string]interface{}{"products": dbResp, "totalCount": count}})
}

func SearchProduct(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	s := c.Query("search")

	pageInt, _ := helper.ConverStringIntoInt(page)
	limitInt, _ := helper.ConverStringIntoInt(limit)
	offsetInt, _ := helper.ConverStringIntoInt(offset)

	products, count, err := database.Mgr.SearchProduct(pageInt, limitInt, offsetInt, s, constant.ProductCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": map[string]interface{}{"products": products, "totalCount": count}})
}

func UpdateProduct(c *gin.Context) {
	userEmail, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": constant.NotRegisteredUser})
		return
	}
	user := database.Mgr.GetSingleRecordByEmailForUser(userEmail.(string), constant.UserCollection)
	if user.UserType != constant.AdminUser {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": constant.NotAuthorizedUserError})
		return
	}
	var updatedReq types.UpdateProduct
	var req types.Product
	err := c.BindJSON(&updatedReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}
	objId, err := primitive.ObjectIDFromHex(updatedReq.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}
	productResp, err := database.Mgr.GetSingleProductById(objId, constant.ProductCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}
	if productResp.Name == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": constant.NoProductAvaliable})
		return
	}
	// old data
	req.ID = objId
	req.Description = productResp.Description
	req.Name = productResp.Name
	req.Price = productResp.Price
	req.MetaInfo = productResp.MetaInfo
	req.ImageUrl = productResp.ImageUrl
	req.CreatedAt = productResp.CreatedAt
	req.UpdatedAt = time.Now().Unix()

	if updatedReq.Name != "" {
		req.Name = updatedReq.Name
	}

	if updatedReq.Description != "" {
		req.Description = updatedReq.Description
	}

	if updatedReq.Price > 0 {
		req.Price = updatedReq.Price
	}
	err = database.Mgr.UpdateProduct(req, constant.ProductCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": updatedReq})
}

func DeleteProduct(c *gin.Context) {
	userEmail, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": constant.NotRegisteredUser})
		return
	}
	user := database.Mgr.GetSingleRecordByEmailForUser(userEmail.(string), constant.UserCollection)
	if user.UserType != constant.AdminUser {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": constant.NotAuthorizedUserError})
		return
	}
	id := c.Query("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}

	productResp, err := database.Mgr.GetSingleProductById(objId, constant.ProductCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}
	if productResp.Name == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": constant.NoProductAvaliable})
		return
	}

	err = database.Mgr.DeleteProduct(objId, constant.ProductCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
}

func CheckoutOrder(c *gin.Context) {
	cartIdStr := c.Param("id")
	cartId, err := primitive.ObjectIDFromHex(cartIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}

	cart, err := database.Mgr.GetCartObjectById(cartId, constant.CartCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}

	cart.Checkout = true

	err = database.Mgr.UpdateCartToCheckOut(cart, constant.CartCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
}
