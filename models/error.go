// models/error.go
package models

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
