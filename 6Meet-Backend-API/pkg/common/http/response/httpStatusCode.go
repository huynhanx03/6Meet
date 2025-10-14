package response

import (
	"net/http"
)

const (
	// Success codes (20000-29999)
	ErrCodeSuccess   = 20001 // Success
	ErrCodeCreated   = 20002 // Resource created successfully
	ErrCodeUpdated   = 20003 // Resource updated successfully
	ErrCodeDeleted   = 20004 // Resource deleted successfully
	ErrCodeRetrieved = 20005 // Resource retrieved successfully

	// Client error codes (40000-49999)
	ErrCodeParamInvalid     = 40001 // Invalid parameters
	ErrCodeValidationFailed = 40002 // Validation failed
	ErrCodeBadRequest       = 40003 // Bad request

	// Authentication/Authorization errors (41000-41999)
	ErrCodeUnauthorized    = 41001 // Unauthorized
	ErrCodeInvalidToken    = 41002 // Invalid token
	ErrCodeTokenExpired    = 41003 // Token expired
	ErrCodeInvalidPassword = 41004 // Invalid password
	ErrCodeAccountNotFound = 41005 // Account not found

	// Forbidden errors (42000-42999)
	ErrCodeForbidden         = 42001 // Forbidden
	ErrCodeInsufficientPerms = 42002 // Insufficient permissions

	// Not found errors (43000-43999)
	ErrCodeNotFound = 43001 // Resource not found

	// Conflict errors (44000-44999)
	ErrCodeConflict = 44001 // Conflict

	// Server error codes (50000-59999)
	ErrCodeInternalServer = 50001 // Internal server error
	ErrCodeDatabaseError  = 50002 // Database error
	ErrCodeMongoDBError   = 50003 // MongoDB error
	ErrCodeRedisError     = 50004 // Redis error
)

var msg = map[int]string{
	// Success
	ErrCodeSuccess:   "Success",
	ErrCodeCreated:   "Resource created successfully",
	ErrCodeUpdated:   "Resource updated successfully",
	ErrCodeDeleted:   "Resource deleted successfully",
	ErrCodeRetrieved: "Resource retrieved successfully",

	// Client errors
	ErrCodeParamInvalid:     "Invalid parameters",
	ErrCodeValidationFailed: "Validation failed",
	ErrCodeBadRequest:       "Bad request",

	// Authentication/Authorization
	ErrCodeUnauthorized:    "Unauthorized",
	ErrCodeInvalidToken:    "Invalid token",
	ErrCodeTokenExpired:    "Token expired",
	ErrCodeInvalidPassword: "Invalid password",
	ErrCodeAccountNotFound: "Account not found",

	// Forbidden
	ErrCodeForbidden:         "Forbidden",
	ErrCodeInsufficientPerms: "Insufficient permissions",

	// Not found
	ErrCodeNotFound: "Resource not found",

	// Conflict
	ErrCodeConflict: "Conflict",

	// Server errors
	ErrCodeInternalServer: "Internal server error",
	ErrCodeDatabaseError:  "Database error",
	ErrCodeMongoDBError:   "MongoDB error",
	ErrCodeRedisError:     "Redis error",
}

func getHTTPStatusCode(code int) int {
	switch {
	case code >= 20000 && code < 30000:
		return http.StatusOK
	case code >= 40000 && code < 41000:
		return http.StatusBadRequest
	case code >= 41000 && code < 42000:
		return http.StatusUnauthorized
	case code >= 42000 && code < 43000:
		return http.StatusForbidden
	case code >= 43000 && code < 44000:
		return http.StatusNotFound
	case code >= 50000 && code < 60000:
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}
