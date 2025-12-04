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

```go
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
```
# 1. What is Logrus

Logrus is a structured, leveled logging library for Go that aims to be a drop-in upgrade over the stdlib `log` package while giving us:

* **Levels** (trace → panic) so we can filter logs by severity.
* **Structured logging**: attach key/value fields to logs (`user=42`, `request_id=abc`).
* **Formatters**: human-readable text or machine-friendly JSON.
* **Hooks**: send log entries to other services (Sentry, Graylog, Kafka, etc.).
* **Custom logger instances** so different packages/services can have separate configs.

Why use it?

* Better for production & observability than plain `log`.
* Easy to integrate with logging backends (ELK/Loki/etc).
* Familiar API: `logrus.WithField(...).Info("...")`.

Tradeoffs:

* More allocations than zero-allocation loggers (e.g., `zerolog`) — fine for most apps, but consider alternatives for ultra-high throughput.

---

# 2. Installing & using Logrus

Install:

```bash
go get github.com/sirupsen/logrus
```

Minimal usage (global logger):

```go
import log "github.com/sirupsen/logrus"

func main() {
    log.SetLevel(log.InfoLevel)
    log.Info("Server started")
}
```

Better pattern — create and configure a logger instance:

```go
logger := logrus.New()
logger.SetLevel(logrus.InfoLevel)
logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
logger.Out = os.Stdout

logger.WithField("service", "auth").Info("auth service started")
```

Important notes:

* Use `logrus.New()` for library code (don’t change global logger defaults).
* For apps, configuring global logger (`logrus.Set...`) is fine in `main()`.
* Avoid `logrus.Fatal` or `logrus.Panic` inside libraries — they exit/panic the process.

---

# 3. LogLevels in Logrus (in depth)

Logrus supports 7 levels (ordered low → high):

* `TraceLevel` — very fine-grained, verbose debugging (lowest).
* `DebugLevel` — debugging info for developers.
* `InfoLevel` — normal operation messages (startup, requests).
* `WarnLevel` — unusual but non-fatal situations.
* `ErrorLevel` — errors which the program can recover from.
* `FatalLevel` — logs + `os.Exit(1)` (terminates).
* `PanicLevel` — logs + `panic()` (stack trace).

Set log level (global or per logger):

```go
logrus.SetLevel(logrus.DebugLevel)
// or for instance
logger.SetLevel(logrus.InfoLevel)
```

Filtering rule: only messages at the configured level *or higher* are emitted.
Example: `SetLevel(WarnLevel)` emits `Warn`, `Error`, `Fatal`, `Panic` only.

Best practices:

* Development: `Trace` or `Debug`.
* Production: `Info` or `Warn` (avoid `Debug` unless diagnosing).
* Don’t overuse `Fatal`/`Panic`. Use them only for unrecoverable app bootstrap errors.

Environment toggling:

```go
lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")) // "debug", "info", ...
if err == nil {
    logrus.SetLevel(lvl)
}
```

---

# 4. Log messages to multiple options (console, file, remote)

Logrus writes to an `io.Writer` (default `os.Stdout`). To log to multiple destinations, use `io.MultiWriter`.

Example: write to console + file (append mode):

```go
f, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    logrus.Fatalf("open log file: %v", err)
}
mw := io.MultiWriter(os.Stdout, f)
logrus.SetOutput(mw)
```

Important details:

* Use `os.OpenFile(..., os.O_APPEND, ...)` — **do not** use `os.Create()` in production because that truncates the file.
* Do the output setup **once** during application startup (not per request).
* For log rotation use a log rotation helper (e.g., `lumberjack.Logger`) — rotate files to avoid runaway disk use.
* For external sinks (Sentry, Graylog, etc.) use **hooks** (`logrus.AddHook(...)`) to push entries asynchronously.

Example with a hook (conceptual):

```go
logrus.AddHook(myHook) // myHook implements logrus.Hook
```

---

# 5. Format messages in Logrus

Logrus has formatters. The two most common are `TextFormatter` (human readable) and `JSONFormatter` (machine readable).

## TextFormatter (human friendly)

Options:

* `FullTimestamp bool` — include full timestamp.
* `TimestampFormat string` — custom timestamp layout.
* `DisableColors bool` / `ForceColors bool` — color control.
* `DisableQuote bool`, `QuoteEmptyFields bool`.

Example:

```go
logrus.SetFormatter(&logrus.TextFormatter{
    FullTimestamp:   true,
    TimestampFormat: time.RFC3339,
})
```

Sample output:

```
time="2025-12-05T15:04:05Z" level=info msg="Server started" service=auth
```

## JSONFormatter (structured & parsable)

Options:

* `TimestampFormat string`
* `DisableTimestamp bool`
* `PrettyPrint bool` — human readable multiline JSON (not recommended for high volume)
* `FieldMap logrus.FieldMap` — change default key names (e.g., `message` → `msg`).

Example:

```go
logrus.SetFormatter(&logrus.JSONFormatter{
    TimestampFormat: time.RFC3339,
    PrettyPrint:     false,
})
```

Sample output:

```json
{"level":"info","msg":"Server started","service":"auth","time":"2025-12-05T15:04:05Z"}
```

Performance tips:

* JSON is better for ingestion by log pipelines.
* PrettyPrint is convenient for debugging but increases bytes and slows logging — keep it off in production.
* Formatter options may allocate; for extremely high performance, consider zero-alloc loggers.

---

# 6. Logging in JSON format (practical details)

Why JSON:

* Easy to index/search in ELK / Loki / Datadog.
* Fields are machine-readable (no ad hoc parsing).
* Compatible with structured tracing/observability pipelines.

Configuration example:

```go
logger := logrus.New()
logger.SetFormatter(&logrus.JSONFormatter{
    TimestampFormat: time.RFC3339Nano,
})
f, _ := os.OpenFile("app.json.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
logger.SetOutput(io.MultiWriter(os.Stdout, f))
logger.SetLevel(logrus.InfoLevel)
```

Add context fields:

```go
logger.WithFields(logrus.Fields{
    "request_id": "abcd-1234",
    "user_id":    42,
}).Info("Processed request")
```

Result (single JSON object per line — good for streaming):

```json
{"level":"info","msg":"Processed request","request_id":"abcd-1234","user_id":42,"time":"2025-12-05T15:04:05Z"}
```

Best practices:

* Keep JSON one object per line (no pretty print) for streaming ingestion.
* Make sure timestamps are in a consistent format (RFC3339 or RFC3339Nano).
* Avoid logging secrets (API keys, passwords) in any output.

---

# 7. `WithField` vs `WithFields` (and how to use them)

Both functions attach structured fields to the log entry and return a `*logrus.Entry`. The entry can be reused (fields persist for the returned entry), and you can chain calls. Differences are just in ergonomics:

### `WithField` — add a single key/value

```go
logrus.WithField("user", "skyy").Info("User logged in")
```

Equivalent to:

```go
logrus.WithFields(logrus.Fields{"user":"skyy"}).Info("User logged in")
```

### `WithFields` — add multiple key/values at once

```go
logrus.WithFields(logrus.Fields{
    "user":       "skyy",
    "request_id": "req-123",
}).Info("User request")
```

### Reuse an entry with many logs

```go
entry := logrus.WithFields(logrus.Fields{
    "service": "payments",
    "env":     "prod",
})

// use entry to log multiple messages with same fields
entry.Info("starting job")
entry.Warn("slow response")
```

This avoids repeating fields in every call.

### `WithError` — special helper to attach an `error`

```go
err := errors.New("db error")
logrus.WithError(err).Error("failed to fetch user")
```

That produces a field named `error` by default.

### Example showing difference in output

Using `WithField`:

```go
logrus.WithField("user", "skyy").Info("login")
```

Output (text):

```
level=info msg="login" user=skyy
```

Using `WithFields`:

```go
logrus.WithFields(logrus.Fields{"user":"skyy","id":101}).Info("login")
```

Output:

```
level=info msg="login" user=skyy id=101
```

### Best practices for fields

* Use short, consistent field names (`request_id`, `user_id`, `service`).
* Attach request/context fields at the earliest point (middleware) and pass the entry down the call chain.
* Don’t create huge dynamic structures as fields (arrays/maps with many items).
* Avoid logging PII or secrets in fields.

---

# Additional practical patterns & tips

## Configure logging once (in `main()`)

```go
func initLogger() *logrus.Logger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: time.RFC3339Nano,
    })
    f, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
    logger.SetOutput(io.MultiWriter(os.Stdout, f))
    logger.SetLevel(logrus.InfoLevel)
    logger.SetReportCaller(true) // include file/func (costly; use when needed)
    return logger
}
```

## Using with Gin

Make a middleware that attaches an entry per request:

```go
func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        entry := logger.WithFields(logrus.Fields{
            "request_id": c.GetHeader("X-Request-ID"),
            "path": c.Request.URL.Path,
            "method": c.Request.Method,
        })
        c.Set("logger", entry)
        entry.Info("request started")
        c.Next()
        entry.WithField("status", c.Writer.Status()).Info("request completed")
    }
}
```

Then inside handlers:

```go
entry := c.MustGet("logger").(*logrus.Entry)
entry.Info("handling business logic")
```

## Hooks for remote sinks

Implement `logrus.Hook`:

```go
type MyHook struct {}
func (h *MyHook) Levels() []logrus.Level { return logrus.AllLevels }
func (h *MyHook) Fire(entry *logrus.Entry) error {
    // push to remote system (async if possible)
    return nil
}
logrus.AddHook(&MyHook{})
```

## Performance considerations

* `WithFields` creates a `map[string]interface{}` and allocates. For very hot code paths, minimize fields or switch to a lower-allocation logger.
* `SetReportCaller(true)` collects runtime.Caller info — useful but relatively expensive.
* JSON formatting and I/O are the main bottlenecks — batch or async transport to remote systems.

## Don’t log passwords/secrets

Always scrub or avoid logging sensitive data. Use structured fields and a filtering layer if necessary.

---

# Quick reference snippets

### Logger setup (JSON, multiwriter)

```go
logger := logrus.New()
logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
f, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
logger.SetOutput(io.MultiWriter(os.Stdout, f))
logger.SetLevel(logrus.InfoLevel)
```

### Add fields and log

```go
logger.WithFields(logrus.Fields{
    "request_id": "r-123",
    "user_id":    42,
}).Info("user request processed")
```

### Use `WithError`

```go
if err != nil {
    logger.WithError(err).WithField("op", "db.query").Error("query failed")
}
```

---
