package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"path"
)

func DataSourceAuthentication(env string) func(c *gin.Context) {
	var token string
	return func(c *gin.Context) {
		if token == "" {
			fmt.Println("Read File")
			var tokenPath string
			if env == "dev" {
				filePath, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				tokenPath = path.Join(filePath, "dev", "data_source.token")
			} else {
				tokenPath = "/data_source_access/data_source.token"
			}
			tokenByte, err := ioutil.ReadFile(tokenPath)
			if err != nil {
				panic(err)
			}
			token = string(tokenByte)
		}
		c.Set("DATA_SOURCE_TOKEN", token)
		c.Next()
	}
}
