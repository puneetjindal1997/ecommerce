package constant

const (
	APIVersion = "v1"

	BadRequestMessage = "request not fulfilled"

	// schedular constants
	HealthCheckRoute = "/health"
	MDBUri           = "localhost:27017"
	Database         = "ecommerce"
	Sender           = ""

	VerifyEmailRoute = "/verify-email"
	VerifyOtpRoute   = "/verify-otp"
	ResendEmailRoute = "/resend-email"
)

const (
	// time slot for otp validation
	OtpValidation = 60
)

// collections
const (
	VerificationsCollection = "verifications"
)

// messages
const (
	EmailValidationError      = "wrong email passed"
	OtpValidationError        = "wrong otp passed"
	OtpExpiredValidationError = "otp expired"
	AlreadyVerifiedError      = "already verified"
	OptAlreadySentError       = "otp already sent to email"
)
