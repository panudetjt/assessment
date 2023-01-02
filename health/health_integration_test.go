//go:build integration

package health

import (
	"io"
	"net/http"
	"testing"

	"github.com/panudetjt/assessment/util"
	"github.com/stretchr/testify/assert"
)

func TestHealthIntegration(t *testing.T) {
	resp := util.Request(http.MethodGet, util.Uri("health"), nil)
	b, _ := io.ReadAll(resp.Body)

	assert.Nil(t, resp.Error)
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "ok", string(b))
}
