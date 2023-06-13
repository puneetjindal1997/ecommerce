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
}
