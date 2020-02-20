package middleware

import (
	"github.com/gin-gonic/gin"
)

const DEV_DATA_SOURCE_HOST string = "127.0.0.1:8081"
const PRODUCT_DATA_SOURCE_HOST string = "funceasy-data-source-service"
func DataSourceService(env string) func(c *gin.Context) {
	return func(c *gin.Context) {
		var dataSourceHost string
		if env == "dev" {
			dataSourceHost = DEV_DATA_SOURCE_HOST
		} else {
			dataSourceHost = PRODUCT_DATA_SOURCE_HOST
		}
		c.Set("DATA_SOURCE_HOST", dataSourceHost)
		c.Next()
	}
}
