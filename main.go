package main

/*

1. What is logrus.
2. Installing & using logrus.
3. LogLevels in logrus.
4. Log messages to multiple options.
5. Format messages in logrus.
6. Logging in JSON format.
7. LogWithField and LogWithFields in logrus.

*/

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {

    logrus.SetReportCaller(true)

    logrus.SetFormatter(&logrus.JSONFormatter{
        DisableTimestamp: true,
        PrettyPrint: true,
    })

    logrus.SetLevel(logrus.TraceLevel)

    // Create file ONCE (append mode)
    f, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        logrus.Fatalln("Error creating log file: ", err)
    }

    // Log to both console and file
    multi := io.MultiWriter(os.Stdout, f)
    logrus.SetOutput(multi)

    // Now logs will go to both
    logrus.Traceln("Trace 🟢")
    logrus.Debugln("Debug 🟡")
    logrus.Infoln("Info 🟠")

    router := gin.New()
    router.GET("/getData", GetDatahandler)

    router.Run(":8081")
}

func GetDatahandler(ctx *gin.Context) {

    logrus.WithField("handler", "GetData").Info("Inside handler")
    logrus.WithFields(logrus.Fields{
        "method": "GetDatahandler",
        "status": "OK",
    }).Info("Handler execution complete")

    ctx.JSON(http.StatusOK, gin.H{
        "data": "Hello from handler 🟢",
    })
}