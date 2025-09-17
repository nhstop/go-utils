package constants

// Error codes (application-wide unique)
const (
	Empty = 0
	// Generic
	InternalServer = 1000
	InvalidRequest = 1001
	Unauthorized   = 1002
	Forbidden      = 1003

	// Crypto / Security
	FailedToEncrypt   = 2000
	FailedToDecrypt   = 2001
	HashingFailed     = 2002
	TokenGeneration   = 2003
	TokenValidation   = 2004
	FailedToGetAESKey = 2005

	// Database
	DBError = 3000

	// Business Logic
	InvalidCredentials = 4000
	UserAlreadyExists  = 4001
	UserNotFound       = 4002
)
