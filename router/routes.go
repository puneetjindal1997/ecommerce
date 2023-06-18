package router

import (
	"net/http"

	"ecommerce/constant"
	"ecommerce/controller"
)

// health check service
var healthCheckRoutes = Routes{
	Route{"Health check", http.MethodGet, constant.HealthCheckRoute, controller.HealthCheck},
}

var userRoutes = Routes{
	Route{"VerifyEmail", http.MethodPost, constant.VerifyEmailRoute, controller.VerifyEmail},
	Route{"VerifyOtp", http.MethodPost, constant.VerifyOtpRoute, controller.VerifyOtp},
	Route{"ResendEmail", http.MethodPost, constant.ResendEmailRoute, controller.VerifyEmail},

	// Register User
	Route{"Register User", http.MethodPost, constant.UserRegisterRoute, controller.RegisterUser},
	Route{"Login User", http.MethodPost, constant.UserLoginRoute, controller.UserLogin},
}

var productGlobalRoutes = Routes{
	Route{"List Product", http.MethodGet, constant.ListProductRoute, controller.ListProductsController},
	Route{"Search Product", http.MethodGet, constant.SearchProductRoute, controller.SearchProduct},
}

var productRoutes = Routes{
	Route{"Register Product", http.MethodPost, constant.RegisterProductRoute, controller.RegisterProduct},
}
