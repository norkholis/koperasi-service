// file: pkg/utils/response.go
package utils

import "github.com/gin-gonic/gin"

func ResponseSuccess(msg string) gin.H {
	return gin.H{"message": msg}
}

func ResponseError(msg string) gin.H {
	return gin.H{"error": msg}
}
