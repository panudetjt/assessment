package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthorization(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	f := func(a string, c echo.Context) (bool, error) {
		if a == "valid" {
			return true, nil
		}
		return false, nil
	}
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	t.Run("return 400 (BadRequest) when no authorization header", func(t *testing.T) {

		h := Authorization(f)(handler)
		got := h(c).(*echo.HTTPError)

		assert.Error(t, got)
		assert.Equal(t, http.StatusBadRequest, got.Code)
	})

	t.Run("return 400 (BadRequest) when authorization header is invalid", func(t *testing.T) {
		req.Header.Set(echo.HeaderAuthorization, "invalid")

		h := Authorization(func(a string, c echo.Context) (bool, error) {
			return false, errors.New("invalid authorization header")
		})(handler)
		got := h(c).(*echo.HTTPError)

		assert.Error(t, got)
		assert.Equal(t, http.StatusBadRequest, got.Code)
	})

	t.Run("return 401 (Unauthorized) when authorization key is invalid", func(t *testing.T) {
		req.Header.Set(echo.HeaderAuthorization, "invalid")

		h := Authorization(func(a string, c echo.Context) (bool, error) {
			return false, nil
		})(handler)
		got := h(c).(*echo.HTTPError)

		assert.Error(t, got)
		assert.Equal(t, http.StatusUnauthorized, got.Code)
	})

	t.Run("return 200 (OK) when authorization key is valid", func(t *testing.T) {
		req.Header.Set(echo.HeaderAuthorization, "valid")

		h := Authorization(f)(handler)
		got := h(c)

		assert.NoError(t, got)
		assert.Equal(t, http.StatusOK, res.Code)
	})
}
