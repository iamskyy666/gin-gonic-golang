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
