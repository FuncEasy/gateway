package middleware

import (
	CRDClient "github.com/funceasy/gateway/pkg/clientset/versioned"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func GetCRDClient(env string) func(c *gin.Context) {
	return func(c *gin.Context) {
		var cfg *rest.Config
		var clientset *CRDClient.Clientset
		var _clientset *kubernetes.Clientset
		var err error
		if env == "dev" {
			var kubeconfig string
			if home := homedir.HomeDir(); home != "" {
				kubeconfig = filepath.Join(home, ".kube", "config")
			} else {
				kubeconfig = ""
			}
			cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				panic(err)
			}
			clientset, err = CRDClient.NewForConfig(cfg)
			_clientset, err = kubernetes.NewForConfig(cfg)
		} else {
			cfg, err := rest.InClusterConfig()
			if err != nil {
				panic(err)
			}
			clientset, err = CRDClient.NewForConfig(cfg)
			_clientset, err = kubernetes.NewForConfig(cfg)
		}
		c.Set("CRDClientSet", clientset)
		c.Set("OriginClientSet", _clientset)
		c.Next()
	}
}
