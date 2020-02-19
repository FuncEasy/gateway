package middleware

import (
	"github.com/funceasy/gateway/pkg/APIError"
	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			e, ok := err.(*APIError.Error)
			if ok {
				c.AbortWithStatusJSON(e.Code, gin.H{
					"error":   e.Err.Error(),
					"message": e.Msg,
				})
				return
			}
			c.AbortWithStatusJSON(500, gin.H{
				"message": "Server Error",
			})
		}
	}()
	c.Next()
}
