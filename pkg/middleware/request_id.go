package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "request_id"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if client sent a Request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set Context
		c.Set(RequestIDKey, requestID)

		// Set Header
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
