package main

import (
	"flag"
	"fmt"
	"github.com/funceasy/gateway/middleware"
	"github.com/funceasy/gateway/router"
	"github.com/gin-gonic/gin"
)

func main() {
	var env *string
	env = flag.String("env", "dev", "run env")
	flag.Parse()
	fmt.Println(*env)

	proxyHost := "localhost:8001"

	r := gin.Default()

	r.Use(middleware.ErrorHandler)
	r.Use(middleware.GetCRDClient(*env))
	r.Use(middleware.DataSourceAuthentication(*env))
	r.Use(middleware.DataSourceService(*env))

	function := r.Group("/function")
	function.POST("/create/:id", router.CreateFunctionCR)
	function.GET("/:id", router.GetFunctionCR)
	function.PUT("update/:id", router.UpdateFunctionCR)
	function.DELETE("/delete/:id", router.DeleteFunctionCR)
	function.POST("/call/:id", router.FunctionCall(*env, proxyHost))

	dataSource := r.Group("/dataSource")
	dataSource.POST("/create", router.CreateDataSource)
	dataSource.POST("/update", router.UpdateDataSource)
	dataSource.DELETE("/:id", router.DeleteDataSource)
	err := r.Run(":8888")
	if err != nil {
		panic(err)
	}
}
