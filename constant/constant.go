package constant

const (
	APIVersion = "v1"

	BadRequestMessage = "request not fulfilled"

	// schedular constants
	HealthCheckRoute = "/health"
	MDBUri           = "localhost:27017"
	Database         = "ecommerce"
	Sender           = ""

	// email verification routes
	VerifyEmailRoute = "/verify-email"
	VerifyOtpRoute   = "/verify-otp"
	ResendEmailRoute = "/resend-email"

	// user related routes
	UserRegisterRoute = "/user-register"
	UserLoginRoute    = "/login"

	// product routes
	RegisterProductRoute = "/product-register"
	ListProductRoute     = "/list-products"
	SearchProductRoute   = "/search"
	UpdateProductRoute   = "/update-product"
	DeleteProductRoute   = "/delete-product"
	AddToCartRoute       = "/cart"
	AddAddressRoute      = "/address"
	GetSingleUserRoute   = "/user/:id"
	UpdateUser           = "/update-user"
	CheckoutRoute        = "/user/:id"
)

const (
	NormalUser = "user"
	AdminUser  = "admin"
)

const (
	// time slot for otp validation
	OtpValidation = 60
)

// collections
const (
	VerificationsCollection = "verifications"
	UserCollection          = "user"
	ProductCollection       = "products"
	AddressCollection       = "user_addresses"
	CartCollection          = "user_cart"
)

// messages
const (
	AlreadyRegisterWithThisEmail = "already register with this email"
	EmailIsNotVerified           = "your email is not verified please verify your email"
	EmailValidationError         = "wrong email passed"
	OtpValidationError           = "wrong otp passed"
	OtpExpiredValidationError    = "otp expired"
	AlreadyVerifiedError         = "already verified"
	OptAlreadySentError          = "otp already sent to email"
	NotRegisteredUser            = "you are not register user"
	PasswordNotMatchedError      = "password doesn't match"
	NotAuthorizedUserError       = "you are not authorized to do this"
	NoProductAvaliable           = "no product avaliable"
	UserDoesNotExists            = "user not exists"
	AddressNotExists             = "address not exists. please add one address"
)
