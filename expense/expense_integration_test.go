//go:build integration

package expense

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	assert.Equal(t, "strawberry smoothie", e.Title)
	assert.Equal(t, 79, e.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", e.Note)
	assert.Equal(t, []string{"food", "beverage"}, e.Tags)
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
