package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
  // Create a Gin router with default middleware (logger and recovery)
  r := gin.Default()

  // Define a simple GET endpoint
  r.GET("/ping", func(c *gin.Context) {
    // Return JSON response
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
	  "status":http.StatusOK,
    })
  })

  // Define a GET endpoint with params
  r.GET("/me/:id", func(c *gin.Context) {
    id:=c.Param("id")
    // Return JSON response
    c.JSON(http.StatusOK, gin.H{
      "user_id":id,
    })
  })

   // Define a simple POST endpoint
  r.POST("/me", func(c *gin.Context) {

    type MeRequest struct{
      Email string `json:"email" binding:"required"`
      Password string `json:"password"`
    }
    var meReq MeRequest

    err:=c.BindJSON(&meReq)
    if err!=nil{
      // Return JSON response
    c.JSON(http.StatusBadRequest, gin.H{
   "error":err.Error(),
   "status_code":http.StatusBadRequest,
    })
    return
    }

    // Return JSON response
    c.JSON(http.StatusOK, gin.H{
   "email":meReq.Email,
   "password":meReq.Password,
    })
  })

  // Define a simple PUT endpoint
  r.PUT("/me", func(c *gin.Context) {

    type MeRequest struct{
      Email string `json:"email" binding:"required"`
      Password string `json:"password"`
    }
    var meReq MeRequest

    err:=c.BindJSON(&meReq)
    if err!=nil{
      // Return JSON response
    c.JSON(http.StatusBadRequest, gin.H{
   "error":err.Error(),
   "status_code":http.StatusBadRequest,
    })
    return
    }

    // Return JSON response
    c.JSON(http.StatusOK, gin.H{
   "email":meReq.Email,
   "password":meReq.Password,
    })
  })

  // Define a simple PATCH endpoint
  r.PATCH("/me", func(c *gin.Context) {

    type MeRequest struct{
      Email string `json:"email" binding:"required"`
      Password string `json:"password"`
    }
    var meReq MeRequest

    err:=c.BindJSON(&meReq)
    if err!=nil{
      // Return JSON response
    c.JSON(http.StatusBadRequest, gin.H{
   "error":err.Error(),
   "status_code":http.StatusBadRequest,
    })
    return
    }

    // Return JSON response
    c.JSON(http.StatusOK, gin.H{
   "email":meReq.Email,
   "password":meReq.Password,
    })
  })

  // Define a simple DELETE endpoint
  r.DELETE("/me/:key", func(c *gin.Context) {

    id:=c.Param("key")


    // Return JSON response
    c.JSON(http.StatusOK, gin.H{
   "id":id,
   "message":"Deleted ✅",
    })
  })

  // Start server on port 8080 (default)
  // Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
  if err := r.Run(); err != nil {
    log.Fatalf("failed to run server: %v", err)
  }
}