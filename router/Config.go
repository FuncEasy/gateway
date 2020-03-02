package router

import (
	"fmt"
	"github.com/funceasy/gateway/pkg/APIError"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetConfig(c *gin.Context) {
	ClientSet, _ := c.Get("OriginClientSet")
	if ClientSet, ok := ClientSet.(*kubernetes.Clientset); ok {
		configMapClient := ClientSet.CoreV1().ConfigMaps("funceasy")
		cm, err := configMapClient.Get("funceasy-config", v1.GetOptions{})
		if err != nil {
			APIError.Panic(err)
		}
		c.JSON(200, gin.H{
			"config": cm.Data,
		})
	} else {
		APIError.PanicError(fmt.Errorf("Failed Get ClientSet "), "Failed Get ClientSet", 500)
	}
}
