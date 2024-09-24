package usecase

import "github.com/gofiber/fiber/v3"

// general error
var (
	ErrFRCreateDatabaseTransaction = fiber.NewError(fiber.StatusInternalServerError, "failed+to+create+database+transaction")
	ErrFRFailedToReadData          = fiber.NewError(fiber.StatusInternalServerError, "failed+to+get+data+from+database")
	ErrFRNotFacebookUser           = fiber.NewError(fiber.StatusConflict, "account+registered+using+Facebook.+Please+login+using+Facebook.")
	ErrFRNotGoogleUser             = fiber.NewError(fiber.StatusConflict, "account+registered+using+Google.+Please+login+using+Google.")
	ErrFRCreateToken               = fiber.NewError(fiber.StatusInternalServerError, "failed+to+create+token")
	ErrFRFailedToStoreData         = fiber.NewError(fiber.StatusInternalServerError, "failed+to+store+data")
	ErrFRFailedToUpdateData        = fiber.NewError(fiber.StatusInternalServerError, "failed+to+update+data")
)

var (
	ErrCreateDatabaseTransaction = fiber.NewError(fiber.StatusInternalServerError, "failed to create database transaction")
	ErrFailedToReadData          = fiber.NewError(fiber.StatusInternalServerError, "failed to get data from database")
	ErrInvalidToken              = fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	ErrFailedToMarshal           = fiber.NewError(fiber.StatusInternalServerError, "failed to marshal")
	ErrFailedToUnMarshal         = fiber.NewError(fiber.StatusInternalServerError, "failed to unmarshal")
	ErrDuplicateUser             = fiber.NewError(fiber.StatusConflict, "duplicate user")
	ErrFailedToClearData         = fiber.NewError(fiber.StatusInternalServerError, "failed to clear data")
	ErrFailedToStoreData         = fiber.NewError(fiber.StatusInternalServerError, "failed to store data")
	ErrFailedToGeneratePassword  = fiber.NewError(fiber.StatusInternalServerError, "failed to generate password")
	ErrFailedToGenerateCode      = fiber.NewError(fiber.StatusInternalServerError, "failed to generate code")
	ErrFailedToSendEmail         = fiber.NewError(fiber.StatusInternalServerError, "failed to send email")
	ErrFailedToVerifyAccount     = fiber.NewError(fiber.StatusConflict, "failed to verify account")
	ErrSendEmailLimitExceeded    = fiber.NewError(fiber.StatusTooManyRequests, "limit to send email exceeded")
	ErrFailedToUpdateData        = fiber.NewError(fiber.StatusInternalServerError, "failed to update data")
	ErrWrongPassword             = fiber.NewError(fiber.StatusUnauthorized, "password you entered is incorrect")
	ErrNeedEmailVerification     = fiber.NewError(fiber.StatusUnauthorized, "please verify your account first")
	ErrLoginTypeOAuth            = fiber.NewError(fiber.StatusUnauthorized, "please login with oauth")
	ErrCreateToken               = fiber.NewError(fiber.StatusInternalServerError, "failed to create token")
	ErrCalculateShipping         = fiber.NewError(fiber.StatusInternalServerError, "failed to calculate shipping cost and estimated days")
	ErrParseStringToNumber       = fiber.NewError(fiber.StatusInternalServerError, "failed to parse string to number")
	ErrOutOfStock                = fiber.NewError(fiber.StatusBadRequest, "failed to add product because out of stock")
)
