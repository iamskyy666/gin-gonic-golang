package main

// What is a middleware
// How to use middleware in Go
// Apply Middleware to routes, routes group and whole application at once

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skyy/gin-gonic/middleware"
)

func main() {
	router := gin.New() // gin-router, without default middleware (New)

	//💡 MW Apply to individual routes
	router.GET("/getData", middleware.Authenticate,middleware.AddHeader,GetDatahandler,) 
	router.GET("/getData1", GetData1handler)
	router.GET("/getData2", GetData2handler)

	//💡 MW Apply to all routes
	//router.Use(middleware.Authenticate) 
	// router.GET("/getData", GetDatahandler)
	// router.GET("/getData1", GetData1handler)
	// router.GET("/getData2", GetData2handler)
	

	// 💡 MW Apply to route-group
	// adminRoutes:=router.Group("/admin",middleware.Authenticate)
	// {
	// adminRoutes.GET("/getData", middleware.Authenticate,GetDatahandler)
	// adminRoutes.GET("/getData1", GetData1handler)
	// adminRoutes.GET("/getData2", GetData2handler)
	// }

	

	// http-config
	server:=&http.Server{
		Addr: ":9091",
		Handler: router,
		ReadTimeout: 10*time.Second,
		WriteTimeout: 10*time.Second,
	}

	err:=server.ListenAndServe()
	if err != nil {
		log.Fatalf("⚠️failed to run server: %v", err)
	}
}

func GetDatahandler(ctx *gin.Context){
	ctx.JSON(http.StatusOK, gin.H{
		"data":"Hi! I am GetDataHandler method() 🟢",
		"status_code":http.StatusOK,
	})
}

func GetData1handler(ctx *gin.Context){
	ctx.JSON(http.StatusOK, gin.H{
		"data":"Hi! I am GetData1Handler method() 🔵",
		"status_code":http.StatusOK,
	})
}

func GetData2handler(ctx *gin.Context){
	ctx.JSON(http.StatusOK, gin.H{
		"data":"Hi! I am GetData2Handler method() 🟡",
		"status_code":http.StatusOK,
	})
}
