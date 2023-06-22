/*
 *	Author: Puneet
 *	Use to start the client
 */
package router

import (
	"ecommerce/auth"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(*gin.Context)
}
type routes struct {
	router *gin.Engine
}

type Routes []Route

/*
 *	Function for grouping ecommerce health routes
 */
func (r routes) EcommerceHealthCheck(rg *gin.RouterGroup) {
	orderRouteGrouping := rg.Group("/ecommerce")
	orderRouteGrouping.Use(CORSMiddleware())
	for _, route := range healthCheckRoutes {
		switch route.Method {
		case "GET":
			orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
		case "OPTIONS":
			orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
		case "PUT":
			orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
		default:
			orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(200, gin.H{
					"result": "Specify a valid http method with this route.",
				})
			})
		}
	}
}

/*
 *	Function for grouping user
 */
func (r routes) EcommerceUser(rg *gin.RouterGroup) {
	orderRouteGrouping := rg.Group("/ecommerce")
	orderRouteGrouping.Use(CORSMiddleware())
	for _, route := range userRoutes {
		switch route.Method {
		case "GET":
			orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
		case "OPTIONS":
			orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
		case "PUT":
			orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
		default:
			orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(200, gin.H{
					"result": "Specify a valid http method with this route.",
				})
			})
		}
	}
}

/*
 *	Function for grouping product
 */
func (r routes) EcommerceProduct(rg *gin.RouterGroup) {
	orderRouteGrouping := rg.Group("/ecommerce")
	orderRouteGrouping.Use(CORSMiddleware())
	for _, route := range productRoutes {
		switch route.Method {
		case "GET":
			orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
		case "OPTIONS":
			orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
		case "PUT":
			orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
		default:
			orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(200, gin.H{
					"result": "Specify a valid http method with this route.",
				})
			})
		}
	}
}

/*
 *	Function for grouping product global routes
 */
func (r routes) EcommerceGlobalProductRoutes(rg *gin.RouterGroup) {
	orderRouteGrouping := rg.Group("/ecommerce-product")
	orderRouteGrouping.Use(CORSMiddleware())
	for _, route := range productGlobalRoutes {
		switch route.Method {
		case "GET":
			orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
		case "OPTIONS":
			orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
		case "PUT":
			orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
		default:
			orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(200, gin.H{
					"result": "Specify a valid http method with this route.",
				})
			})
		}
	}
}

/*
 *	Function for grouping user authentic routes
 */
func (r routes) EcommerceAuthUser(rg *gin.RouterGroup) {
	orderRouteGrouping := rg.Group("/ecommerce")
	orderRouteGrouping.Use(CORSMiddleware())
	for _, route := range userAuthRoutes {
		switch route.Method {
		case "GET":
			orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
		case "OPTIONS":
			orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
		case "PUT":
			orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
		default:
			orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(200, gin.H{
					"result": "Specify a valid http method with this route.",
				})
			})
		}
	}
}

// append routes with versions
func ClientRoutes() {
	r := routes{
		router: gin.Default(),
	}
	v1 := r.router.Group(os.Getenv("API_VERSION"))
	r.EcommerceHealthCheck(v1)
	r.EcommerceUser(v1)
	r.EcommerceGlobalProductRoutes(v1)
	// auth
	v1.Use(auth.Auth())
	r.EcommerceProduct(v1)
	r.EcommerceAuthUser(v1)

	if err := r.router.Run(":" + os.Getenv("PORT")); err != nil {
		log.Println("Failed to run server: %v", err)
	}
}

// Middlewares
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
		}
	}
}
