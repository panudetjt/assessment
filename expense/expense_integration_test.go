//go:build integration

package expense

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/panudetjt/assessment/util"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)
	var e Expense
	res := util.Request(http.MethodPost, util.Uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "strawberry smoothie", e.Title)
	assert.Equal(t, 79, e.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", e.Note)
	assert.Equal(t, []string{"food", "beverage"}, e.Tags)
}

func TestGetById(t *testing.T) {
	e := seedExpense(t)
	var got Expense
	res := util.Request(http.MethodGet, util.Uri("expenses/"+fmt.Sprint(e.ID)), nil)
	err := res.Decode(&got)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.ID, got.ID)
	assert.Equal(t, "strawberry smoothie", got.Title)
	assert.Equal(t, 79, got.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", got.Note)
	assert.Equal(t, []string{"food", "beverage"}, got.Tags)
}

func TestAll(t *testing.T) {
	e := seedExpense(t)

	t.Run("return 401 (Unauthorized) when Authorization key is invalid ", func(t *testing.T) {
		at := os.Getenv("AUTH_TOKEN")
		os.Setenv("AUTH_TOKEN", "invalid")

		res := util.Request(http.MethodGet, util.Uri("expenses"), nil)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// restore
		os.Setenv("AUTH_TOKEN", at)
	})
	t.Run("return all expenses", func(t *testing.T) {
		var got []Expense
		res := util.Request(http.MethodGet, util.Uri("expenses"), nil)
		err := res.Decode(&got)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.GreaterOrEqual(t, len(got), 1)
		assert.Contains(t, got, e)
	})
}

func TestUpdate(t *testing.T) {
	e := seedExpense(t)
	b, _ := json.Marshal(Expense{
		ID:     e.ID,
		Title:  "apple smoothie",
		Amount: 89,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	})

	var got Expense
	res := util.Request(http.MethodPut, util.Uri("expenses/"+fmt.Sprint(e.ID)), strings.NewReader(string(b)))
	err := res.Decode(&got)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.ID, got.ID)
	assert.Equal(t, "apple smoothie", got.Title)
	assert.Equal(t, 89, got.Amount)
	assert.Equal(t, "no discount", got.Note)
	assert.Equal(t, []string{"beverage"}, got.Tags)
}

func seedExpense(t *testing.T) Expense {
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)
	var e Expense
	res := util.Request(http.MethodPost, util.Uri("expenses"), body)
	err := res.Decode(&e)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}
	return e
}
