# Echo Helper
Some useful functions to integrate cleverreach functionality into Echo framework.

## Middlewares
Middlewares can enrich echo or a single handler by complex functionality.<br>
**Make sure you have the according config set up and crtoken initialized!**<br>
Please see [crtoken package](../crtoken/README.md) for more details.

### CheckScope
If you want a jwt token to be checked for a certain scope, use this.<br>
Of course the token is validated first so in your handler you can be sure, everything is fine.
```go
func main() {
    e := echo.New()

    // automaticly validates token and checks for given scope
	e.GET("/dishes", getDishes(), helper.CheckScope("kitchen"))

    e.Start(":8080")
}
```

### GetToken
If you need the crtoken in your handler(s), use this.<br>
You can optionally let it be checked for certain scopes you demand.<br>
Of course the token is validated first so in your handler you can be sure, everything is fine.
```go
func main() {
    e := echo.New()

    // automaticly validates token and stores it in context
    e.GET("/client", func(c echo.Context) error {
        token := c.Get("token").(*crtoken.CRToken)
        c.String("You are client "+token.ClientID)
    }, helper.GetToken())

    // automaticly validates token and stores it in context, if given tokens match
    e.GET("/clientlogin", func(c echo.Context) error {
        token := c.Get("token").(*crtoken.CRToken)
        c.String("client " + token.ClientID + " has login " + token.Login)
    }, helper.GetToken("client", "login"))

    e.Start(":8080")
}
```

## ErrorHandler
Set this as the default error handler of echo for a uniformed error response.

```go
func main() {
    // setup our default error handler for a decent output format
    e.HTTPErrorHandler = helper.EchoErrorHandler
    // ...
}
```

## TestRouter
To be able to easily test your echo handlers, you can use this helper.

```go
// Call your code
func TestPing(t *testing.T) {
	test := assert.New(t)

    // setup test echo
    router := helper.NewTestRouter()

    // put the actual engine into where it's needed
    cmd := &api.Command{Echo: router.Echo}
    cmd.Init()

    // plan a request to one of the handlers...
	router.Request(http.MethodGet, "/ping", nil)
    // ... and fire it. You can do these two lines multiple times on the same router.
	resp := router.Start()

	test.Equal(http.StatusOK, router.Code)
	test.Equal("pong", resp.String())
}
```
Or if you want to use an existing Echo instance:
```go
// Call your code
func TestPing(t *testing.T) {
	test := assert.New(t)

    cmd := &api.Command{}
    cmd.Init()

    // setup test Echo, feeding it the one from where it is existing.
    router := &helper.TestRouter{ Echo: cmd.Echo }

    // plan a request to one of the handlers...
	router.Request(http.MethodGet, "/ping", nil)
    // ... and fire it. You can do these two lines multiple times on the same router.
	resp := router.Start()

	test.Equal(http.StatusOK, router.Code)
	test.Equal("pong", resp.String())
}
```


On a POST request the content type is guessed, except you provide a struct or map, which is marshaled and set to JSON.

If you expect to get called by your code in a test and are using rest.Client{}, [refer to its README.md](../rest/README.md).
