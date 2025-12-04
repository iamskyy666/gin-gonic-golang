package middlewares

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// logger mw

//💡 Write logs to files in GIN.
func FormatLogs(param gin.LogFormatterParams)string{
	return fmt.Sprintf("{%s - [%s] \"%s %s %s %d %s \"%s\" %s\"} \n",
	param.ClientIP,
	param.TimeStamp.Format(time.RFC1123),
	param.Method,
	param.Path,
	param.Request.Proto,
	param.StatusCode,
	param.Latency,
	param.Request.UserAgent(),
	param.ErrorMessage,
)
}

//💡 Logging in JSON format in GIN. (Real world situation).
type logFormatLocal struct{
	TimeStamp time.Time
	StatusCode int
	ClientIP string
	Method string
	Path string
	Latency time.Duration
	RequestProto string
	ErrorMessage string
}


func FormatLogsJSON(param gin.LogFormatterParams)string{
	params:= &logFormatLocal{
	TimeStamp: param.TimeStamp,
	StatusCode: param.StatusCode,
	ClientIP: 	param.ClientIP,
	Method: param.Method,
	Path: param.Path,
	Latency: param.Latency,
	RequestProto: param.Request.Proto,
	ErrorMessage: 	param.ErrorMessage,
	}

	j,err:=json.Marshal(params)
	if err != nil {
		fmt.Println("⚠️failed to marshal! ---", err)
		return err.Error()
	}
	fmt.Println(string(j))
	return  string(j)
	
}