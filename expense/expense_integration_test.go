//go:build integration

package expense

import (
	"bytes"
	"net/http"
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
