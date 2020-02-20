package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/funceasy/gateway/pkg/APIError"
	v1 "github.com/funceasy/gateway/pkg/apis/funceasy.com/v1"
	funceasy "github.com/funceasy/gateway/pkg/clientset/versioned"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"net/http"
)

func CreateFunctionCR(c *gin.Context) {
	var FunctionCRSpec v1.FunctionSpec
	id := c.Param("id")
	if id == "" {
		APIError.PanicError(nil, "ID is Required", 422)
	}
	err := c.ShouldBind(&FunctionCRSpec)
	if err != nil {
		APIError.PanicError(err, "Invalid Input", 422)
	} else {
		if FunctionCRSpec.DataSource != "" {
			token := getDataSourceToken(c, FunctionCRSpec.DataSource)
			FunctionCRSpec.DataServiceToken = token
		}
		function := &v1.Function{
			ObjectMeta: metav1.ObjectMeta{
				Name:      id,
				Namespace: "funceasy",
				Labels: map[string]string{
					"app":      "funceasy_function",
					"function": FunctionCRSpec.Identifier,
				},
			},
			Spec: FunctionCRSpec,
		}
		CRDClientSet, _ := c.Get("CRDClientSet")
		if CRDClientSet, ok := CRDClientSet.(*funceasy.Clientset); ok {
			FunctionClient := CRDClientSet.FunceasyV1().Functions("funceasy")
			_, err = FunctionClient.Create(function)
			if err != nil {
				APIError.Panic(err)
			}
			c.JSON(200, gin.H{
				"message": "success",
			})
		} else {
			APIError.PanicError(fmt.Errorf("Failed Get ClientSet "), "Failed Get ClientSet", 500)
		}
	}
}

func GetFunctionCR(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		APIError.PanicError(nil, "ID is Required", 422)
	}
	CRDClientSet, _ := c.Get("CRDClientSet")
	if CRDClientSet, ok := CRDClientSet.(*funceasy.Clientset); ok {
		FunctionClient := CRDClientSet.FunceasyV1().Functions("funceasy")
		function, err := FunctionClient.Get(id, metav1.GetOptions{})
		if err != nil {
			APIError.PanicError(err, "Not Found", 404)
		}
		c.JSON(200, gin.H{
			"data":    function,
			"message": "success",
		})
	} else {
		APIError.PanicError(fmt.Errorf("Failed Get ClientSet "), "Failed Get ClientSet", 500)
	}
}

func UpdateFunctionCR(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		APIError.PanicError(nil, "ID is Required", 422)
	}
	CRDClientSet, _ := c.Get("CRDClientSet")
	if CRDClientSet, ok := CRDClientSet.(*funceasy.Clientset); ok {
		patch, err := c.GetRawData()
		if err != nil {
			APIError.Panic(err)
		}
		var FunctionCR v1.Function
		err = json.Unmarshal(patch, &FunctionCR)
		if err != nil {
			APIError.Panic(err)
		}
		if FunctionCR.Spec.DataSource != "" {
			token := getDataSourceToken(c, FunctionCR.Spec.DataSource)
			FunctionCR.Spec.DataServiceToken = token
		}
		patch, err = json.Marshal(FunctionCR)
		if err != nil {
			APIError.Panic(err)
		}
		FunctionClient := CRDClientSet.FunceasyV1().Functions("funceasy")
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			function, err := FunctionClient.Get(id, metav1.GetOptions{})
			if err != nil {
				panic(err)
			}
			_, err = FunctionClient.Patch(function.Name, types.MergePatchType, patch)
			return err
		})
		if retryErr != nil {
			APIError.Panic(retryErr)
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	} else {
		APIError.PanicError(fmt.Errorf("Failed Get ClientSet "), "Failed Get ClientSet", 500)
	}
}

func DeleteFunctionCR(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		APIError.PanicError(nil, "ID is Required", 422)
	}
	CRDClientSet, _ := c.Get("CRDClientSet")
	if CRDClientSet, ok := CRDClientSet.(*funceasy.Clientset); ok {
		FunctionClient := CRDClientSet.FunceasyV1().Functions("funceasy")
		err := FunctionClient.Delete(id, &metav1.DeleteOptions{})
		if err != nil {
			APIError.Panic(err)
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	} else {
		APIError.PanicError(fmt.Errorf("Failed Get ClientSet "), "Failed Get ClientSet", 500)
	}
}

func FunctionCall(env string, proxyHost string) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			APIError.PanicError(nil, "ID is Required", 422)
		}
		data, err := c.GetRawData()
		if err != nil {
			APIError.Panic(err)
		}
		var url string
		if env == "dev" {
			url = fmt.Sprintf("http://%s/api/v1/namespaces/funceasy/services/http:function-%s:80/proxy/", proxyHost, id)
		} else {
			url = fmt.Sprintf("http://%s", id)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			APIError.Panic(err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			APIError.Panic(err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			APIError.Panic(err)
		}
		if res.StatusCode != 200 {
			APIError.PanicError(fmt.Errorf(string(body)), "Call Failed", res.StatusCode)
		}
		c.JSON(res.StatusCode, gin.H{
			"res": string(body),
		})
	}
}

func getDataSourceToken(c *gin.Context, dataSourceId string) string {
	DataSourceToken, _ := c.Get("DATA_SOURCE_TOKEN")
	DataSourceServiceHost, _ := c.Get("DATA_SOURCE_HOST")
	data, err := c.GetRawData()
	if err != nil {
		APIError.Panic(err)
	}
	var url = fmt.Sprintf("http://%s/dataSource/token/%s", DataSourceServiceHost, dataSourceId)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
	if err != nil {
		APIError.Panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if DataSourceToken, ok := DataSourceToken.(string); ok {
		req.Header.Set("Authentication", DataSourceToken)
	} else {
		APIError.PanicError(fmt.Errorf("Parse Token to String Failed "), "Parse Token to String Failed", 500)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		APIError.Panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		APIError.Panic(err)
	}
	if res.StatusCode != 200 {
		APIError.PanicError(fmt.Errorf(string(body)), "Get Data Source Token Failed", res.StatusCode)
	}
	fmt.Println("GET TOKEN: ", string(body))
	return string(body)
}
