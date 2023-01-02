package expense

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/panudetjt/assessment/util"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpenseHandler(t *testing.T) {
	t.Run("should return 201 (Created) when request body is valid", func(t *testing.T) {
		e := Expense{
			Title:  "strawberry smoothie",
			Amount: 79,
			Note:   "night market promotion discount 10 bath",
			Tags:   []string{"food", "beverage"},
		}

		b, _ := json.Marshal(e)
		res := util.RequestE(http.MethodPost, "/expenses", strings.NewReader(string(b)))

		e.ID = 12
		db := &util.MockDB{LastInsertId: 12, RowsAffected: 1}
		handler := ExpenseHandler{DB: db}

		handler.CreateExpensesHandler(res.Context)
		var ee Expense
		res.Decode(&ee)

		assert.Equal(t, http.StatusCreated, res.Recorder.Code)
		assert.Equal(t, "INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id", db.Query)
		assert.Equal(t, e, ee)
	})

	t.Run("should return 400 (BadRequest) when request body is invalid", func(t *testing.T) {
		res := util.RequestE(http.MethodPost, "/expenses", strings.NewReader("invalid body"))
		db := &util.MockDB{LastInsertId: 12, RowsAffected: 1}
		handler := ExpenseHandler{DB: db}

		handler.CreateExpensesHandler(res.Context)
		var e util.Error
		res.Decode(&e)

		assert.Equal(t, http.StatusBadRequest, res.Recorder.Code)
		assert.NotEmpty(t, e.Message)
	})
	t.Run("should return 500 (InternalServerError) when database error", func(t *testing.T) {
		e := Expense{
			Title:  "strawberry smoothie",
			Amount: 79,
			Note:   "night market promotion discount 10 bath",
			Tags:   []string{"food", "beverage"},
		}

		b, _ := json.Marshal(e)
		res := util.RequestE(http.MethodPost, "/expenses", strings.NewReader(string(b)))

		row := &util.MockRowError{MockRow: util.MockRow{LastInsertId: 0, RowsAffected: 0}}
		db := &util.MockDB{LastInsertId: 0, RowsAffected: 0, Row: row}
		handler := ExpenseHandler{DB: db}

		handler.CreateExpensesHandler(res.Context)
		var err util.Error
		res.Decode(&err)

		assert.Equal(t, http.StatusInternalServerError, res.Recorder.Code)
		assert.NotEmpty(t, err.Message)
	})
}
