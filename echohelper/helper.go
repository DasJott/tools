package helper

import (
	"fmt"
	"net/http"

	"cleverreach.com/crtools/crtoken"
	"cleverreach.com/crtools/rest"
	"github.com/labstack/echo/v4"
)

// EchoErrorHandler can be set as an error handler in an echo object
func EchoErrorHandler(err error, c echo.Context) {
	c.JSON(rest.Error(http.StatusInternalServerError, err.Error()))
}

// CheckScope is a middleware func for echo to check a certain scope
func CheckScope(scope string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tok, err := crtoken.FromRequest(c.Request())
			if err != nil {
				fmt.Println(err.Error())
				return c.JSON(rest.Error(http.StatusBadRequest, "invalid token"))
			}
			if tok.HasScope(scope) {
				return next(c)
			}
			return c.JSON(rest.Error(http.StatusForbidden, "scope mismatch"))
		}
	}
}

// GetToken retrieves the token with given scopes.
// You get it from Context by `c.Get("token").(*crtoken.CRToken)`
func GetToken(scopes ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := crtoken.FromRequest(c.Request())
			if err != nil {
				fmt.Println(err.Error())
				return c.JSON(rest.Error(http.StatusBadRequest, "invalid token"))
			}
			if scopes == nil || token.MatchScopes(scopes) {
				c.Set("token", token)
				return next(c)
			}
			return c.JSON(rest.Error(http.StatusForbidden, "scope mismatch"))
		}
	}

}
