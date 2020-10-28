# Rest
This package contains useful functionality for REST APIs.<br>

## Client
A easy to use rest client to send JSON information around.<br>
You can use it right away or might want to initialize it.

### Environment / crconfig variables
This package is somewhat configurable. SOme values can be set on the package, some by config/environment.<br>
Environment wins over config wins over package settings. [See here for more details on crconfig](../crconfig/README.md)

Here are the values to be set:
- REST_CLIENT_TIMEOUT<br>
A connection timeout. It also includes the reading from response, so be careful.<br>
You can set this value individually for each call using the Timeout() function.<br>
Default is 0 (which is no timeout).
- REST_CLIENT_RETRIES<br>
Number of retries to make after response status >= 500.<br>
To set this individually for a request, use Retries() function.<br>
Default is 0.
- REST_CLIENT_MAX_RETRY_TIME<br>
Maximum time in seconds to use to wait between retries.<br>
To set this individually for a request, use Retries() function.<br>
Default is 60s.
- REST_CLIENT_IGNORE_CERTIFICATE<br>
Flag to make the client ignore any certifcate errors.<br>
Default is false.
- REST_CLIENT_BASE_PATH (rest.BasePath)<br>
The base path to be used and prefixed on every call to NewClient().<br>
- DEBUG (rest.Debug)<br>
Switch on some debug functionality. Mostly logging.<br>

### Usage
Create an instance providing the URL.<br>
There are two ways to get an instance:
```go
// working with a set BasePath
client := rest.NewClient("family", "members", 2)

// defining complete URL here
client := rest.NewClient("https://rest.myapi.com", "family", "members", 2)
```
Or if you have a complete URL from somewhere:
```go
// completeURL := "https://rest.myapi.com/family/members?gender=female&parent=true"
client := rest.NewClientFromURL(completeURL)
```

### Example
```go
// dayNum := 3
day := rest.NewClient("https://rest.week.com/days", dayNum).Get().String()
// day == "Wednesday"
```
There are lots of methods to send the request, please refer to suggestions of your Editor.

### CatchRequest (for testing purposes)
If you want to test code using rest.Client, you can easily catch those requests to reply as your test wants to.<br>
Simply provide a function to field `rest.CatchRequest`. Example follows:
```go
func TestMyCall(t *testing.T) {
    rest.CatchRequest = func(req *http.Request) *http.Response {
        rec := httptest.NewRecorder()
        rec.WriteString("boo!")
        return rec.Result()
    }
    // ...
}
```

### Logger
You can set a certain logger for retrieving debug information.<br>
The logger must implement following interface, which should match most loggers.
```go
type ClientLogger interface {
    Warnln(...interface{})
    Debugln(...interface{})
}
```

## Error
This function fits most JSON response functions of frameworks like echo or gin.<br>
They require code and interface{} as parameters and this is exactly what this function returns.

### Example
```go
// for echo
func handleCall(c echo.Context) error {
    return c.JSON(rest.Error(http.StatusBadRequest, "error creating token"))
}

// same for gin
func handleCall(c *gin.Context) error {
    return c.JSON(rest.Error(http.StatusBadRequest, "error creating token"))
}
```