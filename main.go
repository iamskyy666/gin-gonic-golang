package main

import (
	"log"
	"net/http"

	"io"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default() // gin-router, with default middleware

	router.GET("/", RootHandler)
	router.POST("/", PostHandler)
	router.GET("/get-body-data", GetBodyDataHandler)
	router.GET("/get-QryStr", GetQryDataHandler)
	router.GET("/get-UrlParams/:name/:age", GetUrlDataHandler)

	err := router.Run() //default/without params:8080
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

// ROOT
func RootHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"Thought💭": "Don't take life too seriously, you ain't getting out alive anyways",
		"data":     "🍸 Welcome to GET root / home router Skyy (PORT: 8️⃣0️⃣8️⃣0️⃣ by default)!",
		"status":   http.StatusOK,
	})
}

// POST
func PostHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"data":   "Hi I'm a POST request 🟢",
		"status": http.StatusOK,
	})
}

// GET
func GetBodyDataHandler(ctx *gin.Context) {
  // Read data from the body
	body := ctx.Request.Body
	val, err := io.ReadAll(body)

  if err!=nil{
    ctx.JSON(http.StatusInternalServerError, gin.H{
		"ERROR ⚠️": err.Error(),
		"status":   http.StatusInternalServerError,
	})
  log.Fatal(err.Error())
  return
  }

	ctx.JSON(http.StatusOK, gin.H{
		"bodyData": string(val),
		"status":   http.StatusOK,
	})
}

// Handling query-params
// http://localhost:8080/get-QryStr?name=Mark&age=30
// GET
func GetQryDataHandler(ctx *gin.Context) {
  // Read data from the body
	name := ctx.Query("name")
  age := ctx.Query("age")

	ctx.JSON(http.StatusOK, gin.H{
		"data": "Getting data from Query-Params 🟢",
    "name":name,
    "age":age,
		"status":   http.StatusOK,
	})
}

// Handling URL-params
// http://localhost:8080/get-UrlParams/Skyy/30
// GET
func GetUrlDataHandler(ctx *gin.Context) {
  // Read data from the URL-params
	name := ctx.Param("name")
  age := ctx.Param("age")

	ctx.JSON(http.StatusOK, gin.H{
		"data": "Getting data from URL params 🔵",
    "name":name,
    "age":age,
		"status":   http.StatusOK,
	})
}
