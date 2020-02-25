package main

import (
	"flag"
	"fmt"
	"github.com/funceasy/gateway/middleware"
	"github.com/funceasy/gateway/router"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	var env *string
	env = flag.String("env", "dev", "run env")
	flag.Parse()
	fmt.Println(*env)
	port := os.Getenv("GATEWAY_SERVICE_PORT")
	if port == "" {
		port = ":8082"
	} else {
		port = ":" + port
	}
	proxyHost := "localhost:8001"

	r := gin.Default()

	r.Use(middleware.ErrorHandler)
	r.Use(middleware.Authentication(*env))
	r.Use(middleware.GetCRDClient(*env))
	r.Use(middleware.DataSourceAuthentication(*env))
	r.Use(middleware.DataSourceService(*env))

	function := r.Group("/function")
	function.POST("/call/:id", router.FunctionCall(*env, proxyHost, "POST"))
	function.GET("/call/:id", router.FunctionCall(*env, proxyHost, "GET"))
	function.GET("/instance/:id", router.GetFunctionCR)
	function.POST("/create/:id", router.CreateFunctionCR)
	function.PUT("update/:id", router.UpdateFunctionCR)
	function.DELETE("/delete/:id", router.DeleteFunctionCR)

	dataSource := r.Group("/dataSource")
	dataSource.POST("/create", router.CreateDataSource)
	dataSource.POST("/update", router.UpdateDataSource)
	dataSource.DELETE("/:id", router.DeleteDataSource)
	err := r.Run(port)
	if err != nil {
		panic(err)
	}
}
