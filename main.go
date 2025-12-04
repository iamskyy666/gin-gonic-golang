package main

/*

1. Logging in GIN.
2. How default logging works.
3. Define format for the log of routes in GIN.
4. Define format of the logs with GIN.
5. Write logs to files in GIN.
6. Controlling log-output coloring in console with GIN.
7. Logging in JSON format in GIN. (Real world situation).

*/

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"github.com/skyy/gin-gonic/middlewares"
)

func main() {
	// router := gin.Default()
	router := gin.Default()

	// Define format for the log of routes
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	log.Printf(
		"Custom Route Log → method=%s | path=%s | handler=%s | handlers=%d",
		httpMethod, absolutePath, handlerName, nuHandlers,
	)
}

  // Controlling log-output coloring in console with GIN. 🎨
  gin.ForceConsoleColor()
  gin.DefaultWriter = colorable.NewColorableStdout()


	// Create a log-file and write logs (data) to it.
	f,_:=os.Create("ginLogging.log")
	//gin.DefaultWriter = io.MultiWriter(f) // log to file
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) // log to file + console

	//💡 Logger MW
	// router.Use(gin.LoggerWithFormatter(middlewares.FormatLogs))

	// 💡 JSON logger
	router.Use(gin.LoggerWithFormatter(middlewares.FormatLogsJSON))


	router.GET("/getData",GetDatahandler) 

	err:=router.Run()
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