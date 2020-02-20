package router

import (
	"bytes"
	"fmt"
	"github.com/funceasy/gateway/pkg/APIError"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func CreateDataSource(c *gin.Context) {
	DataSourceToken, _ := c.Get("DATA_SOURCE_TOKEN")
	DataSourceServiceHost, _ := c.Get("DATA_SOURCE_HOST")
	data, err := c.GetRawData()
	if err != nil {
		APIError.Panic(err)
	}
	var url = fmt.Sprintf("http://%s/dataSource/create", DataSourceServiceHost)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
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
		APIError.PanicError(fmt.Errorf(string(body)), "Create Failed", res.StatusCode)
	}
	c.JSON(res.StatusCode, gin.H{
		"res": string(body),
	})
}

func UpdateDataSource(c *gin.Context) {
	DataSourceToken, _ := c.Get("DATA_SOURCE_TOKEN")
	DataSourceServiceHost, _ := c.Get("DATA_SOURCE_HOST")
	data, err := c.GetRawData()
	if err != nil {
		APIError.Panic(err)
	}
	var url = fmt.Sprintf("http://%s/dataSource/update", DataSourceServiceHost)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
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
		APIError.PanicError(fmt.Errorf(string(body)), "Update Failed", res.StatusCode)
	}
	c.JSON(res.StatusCode, gin.H{
		"res": string(body),
	})
}

func DeleteDataSource(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		APIError.PanicError(nil, "ID is Required", 422)
	}
	DataSourceToken, _ := c.Get("DATA_SOURCE_TOKEN")
	DataSourceServiceHost, _ := c.Get("DATA_SOURCE_HOST")
	data, err := c.GetRawData()
	if err != nil {
		APIError.Panic(err)
	}
	var url = fmt.Sprintf("http://%s/dataSource/%s", DataSourceServiceHost, id)
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(data))
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
		APIError.PanicError(fmt.Errorf(string(body)), "Delete Failed", res.StatusCode)
	}
	c.JSON(res.StatusCode, gin.H{
		"res": string(body),
	})
}
