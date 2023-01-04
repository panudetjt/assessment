package expense

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
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
		db, mock, _ := sqlmock.New()
		mock.ExpectQuery("INSERT INTO expenses (.+) VALUES (.+) RETURNING id").
			WithArgs(e.Title, e.Amount, e.Note, pq.Array(e.Tags)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
		handler := Handler{DB: db}
		e.ID = 1

		handler.CreateExpensesHandler(res.Context)
		var ee Expense
		res.Decode(&ee)

		assert.Equal(t, http.StatusCreated, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Equal(t, e, ee)
	})

	t.Run("should return 400 (BadRequest) when request body is invalid", func(t *testing.T) {
		res := util.RequestE(http.MethodPost, "/expenses", strings.NewReader("invalid body"))
		db, mock, _ := sqlmock.New()
		handler := Handler{DB: db}

		handler.CreateExpensesHandler(res.Context)
		var e util.Error
		res.Decode(&e)

		assert.Equal(t, http.StatusBadRequest, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
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
		db, mock, _ := sqlmock.New()
		mock.ExpectQuery("INSERT INTO expenses (.+) VALUES (.+) RETURNING id").
			WithArgs(e.Title, e.Amount, e.Note, pq.Array(e.Tags)).
			WillReturnError(&pq.Error{})
		handler := Handler{DB: db}

		handler.CreateExpensesHandler(res.Context)
		var err util.Error
		res.Decode(&err)

		assert.Equal(t, http.StatusInternalServerError, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotEmpty(t, err.Message)
	})
}
