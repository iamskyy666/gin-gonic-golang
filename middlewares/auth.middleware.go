package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//💡 auth req-middleware
func Authenticate(ctx *gin.Context){
	if !(ctx.Request.Header.Get("Token")=="auth"){
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{
			"Message":"Token Not Present! 🔴",
		})
		return
	}

	ctx.Next()
}

// 💡Alternate way to write same MW
//func Authenticate()gin.HandlerFunc{
	// Write custom logic to be applied before the MW is executed
// 	return func(ctx *gin.Context){
// 	if !(ctx.Request.Header.Get("Token")=="auth"){
// 		ctx.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{
// 			"Message":"Token Not Present! 🔴",
// 		})
// 		return	
// 	}
// 		ctx.Next()
// 	}
// }

// 💡 resp-middleware (runs before the resp. is executed)
func AddHeader(ctx *gin.Context){
	ctx.Writer.Header().Set("Key","Val")
	ctx.Next()
}
