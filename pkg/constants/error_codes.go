package constants

// Error codes (application-wide unique)
const (
	// Generic
	ErrCodeInternalServer = 1000
	ErrCodeInvalidRequest = 1001
	ErrCodeUnauthorized   = 1002
	ErrCodeForbidden      = 1003

	// Crypto / Security
	ErrCodeFailedToEncrypt = 2000
	ErrCodeFailedToDecrypt = 2001
	ErrCodeHashingFailed   = 2002
	ErrCodeTokenGeneration = 2003
	ErrCodeTokenValidation = 2004

	// Database
	ErrCodeCreateUserFailed = 3000
	ErrCodeUpdateUserFailed = 3001
	ErrCodeDeleteUserFailed = 3002
	ErrCodeFetchUserFailed  = 3003

	// Business Logic
	ErrCodeInvalidCredentials = 4000
	ErrCodeUserAlreadyExists  = 4001
	ErrCodeUserNotFound       = 4002
)
