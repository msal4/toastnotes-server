package utils

import "github.com/gin-gonic/gin"

// Msg generates a success response with the provided message.
func Msg(msg string) gin.H {
	return gin.H{"message": msg}
}

// Err generates a failed response with the provided message.
func Err(msg string) gin.H {
	return gin.H{"error": msg}
}
