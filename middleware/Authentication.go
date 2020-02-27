package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/funceasy/gateway/pkg/APIError"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

func Authentication(env string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authentication")
		if tokenStr == "" {
			tokenStr = c.GetHeader("authentication")
		}
		if tokenStr == "" {
			APIError.PanicError(fmt.Errorf("No Authentication. "), "Forbidden", 401)
		} else {
			token, err := jwt.Parse(tokenStr, getPublicKey(env))
			if err != nil {
				APIError.PanicError(fmt.Errorf("Verify Failed. "), "Forbidden", 401)
			}
			if !token.Valid {
				APIError.PanicError(fmt.Errorf("Not Authorized. "), "Forbidden", 401)
			} else {
				c.Next()
			}
		}
	}
}

func getPublicKey(env string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Invaild Token. ")
		} else {
			var publicKeyPath string
			if env == "dev" {
				filePath, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				publicKeyPath = path.Join(filePath, "dev", "gateway.public.key")
			} else {
				publicKeyPath = "/gateway_access/gateway.public.key"
			}
			publicKeyByte, err := ioutil.ReadFile(publicKeyPath)
			if err != nil {
				logrus.Error(err)
				return nil, fmt.Errorf("Read Key Error. ")
			}
			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			return publicKey, nil
		}
	}
}
