```go
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
```
```go
// Custom HTTP Config. with GIN
// Route Grouping in GIN
// Basic Auth funtionality in GIN
func main() {
	router := gin.Default() // gin-router, with default middleware

	// Without group
	 router.GET("/", RootHandler)

	//💡 Auth 🛡️
	auth:=gin.BasicAuth(gin.Accounts{
		"user":"passw",
		"user1":"passw1",
		"user2":"passw2",
	})

	//💡 Grouping routes 🛜
	adminRoutes:= router.Group("/admin",auth) // auth applied
	{
		adminRoutes.GET("/get-body-data", GetBodyDataHandler).GET("/get-QryStr", GetQryDataHandler).GET("/get-UrlParams/:name/:age", GetUrlDataHandler)
	}

	clientRoutes:= router.Group("/client")
	{
		clientRoutes.GET("/get-UrlParams/:name/:age", GetUrlDataHandler)
	}

	//💡 custom http-config ⚙️
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
```
1️⃣ **Why we use `&http.Server{}` instead of `router.Run()`**
2️⃣ **How `gin.BasicAuth()` works**


# ✅ 1. Advantage of using `&http.Server{}` over `router.Run()`

In Gin, we can start a server in two ways:

### **(A) Simple way → `router.Run()`**

```go
router.Run(":9091")
```

This method is **simple and quick**, but limited.

### **(B) Advanced way → custom `http.Server{}`**

```go
server := &http.Server{
    Addr:         ":9091",
    Handler:      router,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
}
server.ListenAndServe()
```

---

## 🎯 **Advantages of `http.Server{}`**

### ✔ 1. **Timeout controls**

With `ReadTimeout`, `WriteTimeout`, `IdleTimeout` etc., we can:

* stop slow clients
* protect server from DDoS-like behavior
* prevent unreadable long requests

`router.Run()` does NOT offer these options.

Example:

```go
ReadTimeout:  10 * time.Second,
WriteTimeout: 10 * time.Second,
```

---

### ✔ 2. **TLS / HTTPS support**

With `http.Server`, we can run:

```go
server.ListenAndServeTLS("cert.pem", "key.pem")
```

`router.Run()` **can’t configure custom TLS**.

---

### ✔ 3. **Graceful shutdown**

We can gracefully stop the server using `server.Shutdown(ctx)`:

* finish ongoing requests
* avoid dropping connections
* useful for microservices & production

`router.Run()` does NOT support graceful shutdown.

---

### ✔ 4. **Custom server configurations**

We can configure:

* Max header size
* Keep-alive settings
* Custom connection state handling
* Logging
* HTTP/2 settings

All these are not possible with `router.Run()`.

---

## 🎉 **Conclusion**

`router.Run()` → good for **simple/testing**
`&http.Server{}` → required for **production**, secure, configurable, scalable servers.

---

# ✅ 2. Understanding `gin.BasicAuth`

`gin.BasicAuth()` is middleware that protects routes using **HTTP Basic Authentication**.

### ✓ BasicAuth stores allowed username–password pairs:

```go
auth := gin.BasicAuth(gin.Accounts{
    "user":  "passw",
    "user1": "passw1",
    "user2": "passw2",
})
```

---

## ⚙️ **How BasicAuth works internally**

1. When a request arrives, Gin checks the HTTP Header:

   ```
   Authorization: Basic base64(username:password)
   ```

2. If the header is missing → Gin returns:

   ```
   401 Unauthorized
   ```

3. If the username/password is wrong → Gin returns:

   ```
   401 Unauthorized
   ```

4. If correct → request passes to the next handler.

---

## 🔐 Example of using BasicAuth with route groups

```go
admin := router.Group("/admin", auth)
{
    admin.GET("/dashboard", DashboardHandler)
}
```

Now only valid users (user/passw, user1/passw1…) can access `/admin/*`.

---

## 🎯 Use Cases of BasicAuth

* Admin dashboards
* Developer-only testing routes
* Local APIs
* Quick security for internal tools

⚠️ **Not recommended for public production APIs**
Use **JWT** or **OAuth** for serious authentication.

---

# ⭐ Final Summary

### **Why use `http.Server{}`?**

| Feature                      | router.Run() | http.Server{}  |
| ---------------------------- | ------------ | -------------- |
| Read/Write timeout           | ❌ No         | ✅ Yes          |
| TLS setup                    | ❌ Limited    | ✅ Full control |
| Graceful shutdown            | ❌ No         | ✅ Yes          |
| Max header/connection config | ❌ No         | ✅ Yes          |
| Production-ready             | ❌ Not really | ✅ Yes          |

### **gin.BasicAuth**

* Simple username/password checking middleware
* Protects routes
* Sends `401` if unauthorized
* Good for local/internal use, not heavy production use

---

```go
package middleware

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
```
```go
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
```

# ✅ **Two Styles of Middleware in Gin**

## ✔ Style 1 — Direct Handler Function

### **(Used in your `Authenticate(ctx *gin.Context)` example)**

```go
func Authenticate(ctx *gin.Context) {
	if !(ctx.Request.Header.Get("Token") == "auth") {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Message": "Token Not Present! 🔴",
		})
		return
	}
	ctx.Next()
}
```

### 🔍 **Characteristics**

* Middleware is written directly as a function with the signature:

  ```
  func(ctx *gin.Context)
  ```
* Simple and direct.
* Cannot accept parameters.
* Only works when directly passed as:

  ```go
  router.Use(Authenticate)
  admin.Use(Authenticate)
  ```

---

## ✔ Style 2 — Middleware Factory (Returning `gin.HandlerFunc`)

### **(Used in your alternate `Authenticate()` example)**

```go
func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Header.Get("Token") != "auth" {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"Message": "Token Not Present! 🔴",
			})
			return
		}
		ctx.Next()
	}
}
```

### 🔍 **Characteristics**

* Function returns another function (`gin.HandlerFunc`).
* More flexible and reusable.
* Allows passing **arguments/configuration** to middleware.

Example:

```go
func AuthenticateWithToken(expected string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Header.Get("Token") != expected {
			ctx.AbortWithStatusJSON(401, gin.H{"Message": "Invalid token!"})
			return
		}
		ctx.Next()
	}
}
```

Usage:

```go
router.Use(AuthenticateWithToken("auth123"))
```

---

# 🎯 **Key Differences (Very Important)**

| Feature                            | Style 1 (Direct)         | Style 2 (Factory)        |
| ---------------------------------- | ------------------------ | ------------------------ |
| Function signature                 | `func(ctx *gin.Context)` | `func() gin.HandlerFunc` |
| Can pass parameters?               | ❌ No                     | ✅ Yes                    |
| Best for simple logic              | ✔ Yes                    | ✔ Yes                    |
| Can create configurable middleware | ❌ No                     | ✔ Yes                    |
| More reusable/flexible             | ❌ No                     | ✔ Yes                    |
| How we use it                      | `.Use(Authenticate)`     | `.Use(Authenticate())`   |

---

# 🧠 **Why Style 2 is more powerful**

If we want middleware with variables, settings, or custom behavior, Style 2 is the only choice.

Example:
Middleware to check for **different tokens**:

```go
admin.Use(AuthenticateWithToken("ADMIN123"))
client.Use(AuthenticateWithToken("CLIENT123"))
```

You cannot do this with the Style 1 middleware.

---

# ⭐ Summary

### **Style 1 (Direct middleware):**

* Simpler
* Good for fixed, static logic
* No dynamic parameters
* Used as: `router.Use(Authenticate)`

### **Style 2 (Middleware factory):**

* More flexible
* Can accept parameters (configurable middleware)
* Used as: `router.Use(Authenticate())`
* Best for real-world production apps

---

```go
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
```
```go
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
```