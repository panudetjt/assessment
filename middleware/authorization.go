package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthorizationValidator func(string, echo.Context) (bool, error)

func Authorization(f AuthorizationValidator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			if auth == "" {
				return echo.NewHTTPError(http.StatusBadRequest, errors.New("missing authorization header"))
			}
			ok, err := f(auth, c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid authorization header"))
			}
			if !ok {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
