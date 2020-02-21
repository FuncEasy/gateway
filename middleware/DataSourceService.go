package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
)

const DEV_DATA_SOURCE_HOST string = "127.0.0.1:8081"
const PRODUCT_DATA_SOURCE_SERVICE string = "funceasy-data-source-service"
func DataSourceService(env string) func(c *gin.Context) {
	return func(c *gin.Context) {
		var dataSourceHost string
		if env == "dev" {
			dataSourceHost = DEV_DATA_SOURCE_HOST
		} else {
			host := os.Getenv("DATA_SOURCE_SERVICE")
			if host != "" {
				dataSourceHost = host
			} else {
				dataSourceHost = PRODUCT_DATA_SOURCE_SERVICE
			}
		}
		c.Set("DATA_SOURCE_HOST", dataSourceHost)
		c.Next()
	}
}
