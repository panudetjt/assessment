package expense

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/panudetjt/assessment/util"
	"github.com/stretchr/testify/assert"
)

func TestUpdateExpenseHandler(t *testing.T) {
	t.Run("should return 200 (OK) when request body is valid", func(t *testing.T) {
		e := Expense{
			ID:     1,
			Title:  "apple smoothie",
			Amount: 89,
			Note:   "no discount",
			Tags:   []string{"beverage"},
		}
		b, _ := json.Marshal(e)
		res := util.RequestE(http.MethodPut, "/expenses/1", strings.NewReader(string(b)))
		res.Context.SetPath("/expenses/:id")
		res.Context.SetParamNames("id")
		res.Context.SetParamValues("1")
		mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow("1", "test-title", "123", "test-note", pq.Array([]string{"test-tags"}))
		db, mock, _ := sqlmock.New()
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = ?").
			ExpectQuery().
			WithArgs("1").
			WillReturnRows(mockRows)
		mock.ExpectPrepare("UPDATE expenses SET title = \\$2, amount = \\$3, note = \\$4, tags = \\$5 WHERE id = \\$1").
			ExpectExec().
			WithArgs(e.ID, e.Title, e.Amount, e.Note, pq.Array(e.Tags)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		handler := Handler{DB: db}

		handler.UpdateExpensesHandler(res.Context)
		ee := Expense{}
		res.Decode(&ee)

		assert.Equal(t, "1", res.Context.Param("id"))
		assert.Equal(t, http.StatusOK, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Equal(t, e, ee)
	})
	t.Run("should return 400 (BadRequest) when request body is invalid", func(t *testing.T) {
		res, db, mock := arrange("invalid body")
		mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow("1", "test-title", "123", "test-note", pq.Array([]string{"test-tags"}))
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = ?").
			ExpectQuery().
			WithArgs("1").
			WillReturnRows(mockRows)

		handler := Handler{DB: db}
		handler.UpdateExpensesHandler(res.Context)
		ee := util.Error{}
		res.Decode(&ee)

		assert.Equal(t, "1", res.Context.Param("id"))
		assert.Equal(t, http.StatusBadRequest, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotNil(t, ee.Message)
	})
	t.Run("should return 404 (NotFound) when not found row", func(t *testing.T) {
		res, db, mock := arrange("")
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = ?").
			ExpectQuery().
			WillReturnError(sql.ErrNoRows)

		handler := Handler{DB: db}
		handler.UpdateExpensesHandler(res.Context)
		ee := util.Error{}
		res.Decode(&ee)

		assert.Equal(t, "1", res.Context.Param("id"))
		assert.Equal(t, http.StatusNotFound, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotNil(t, ee.Message)
	})
	t.Run("should return 500 (InternalServerError) when cannot prepare SELECT", func(t *testing.T) {
		res, db, mock := arrange("")
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = ?").WillReturnError(&pq.Error{})

		handler := Handler{DB: db}
		handler.UpdateExpensesHandler(res.Context)
		ee := util.Error{}
		res.Decode(&ee)

		assert.Equal(t, "1", res.Context.Param("id"))
		assert.Equal(t, http.StatusInternalServerError, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotNil(t, ee.Message)
	})
	t.Run("should return 500 (InternalServerError) when cannot execute SELECT", func(t *testing.T) {
		res, db, mock := arrange("")
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = ?").
			ExpectQuery().
			WillReturnError(&pq.Error{})

		handler := Handler{DB: db}
		handler.UpdateExpensesHandler(res.Context)
		ee := util.Error{}
		res.Decode(&ee)

		assert.Equal(t, "1", res.Context.Param("id"))
		assert.Equal(t, http.StatusInternalServerError, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotNil(t, ee.Message)
	})
	t.Run("should return 500 (InternalServerError) when cannot prepare update", func(t *testing.T) {
		res, db, mock := arrange("")
		mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow("1", "test-title", "123", "test-note", pq.Array([]string{"test-tags"}))
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = ?").
			ExpectQuery().
			WithArgs("1").
			WillReturnRows(mockRows)
		mock.ExpectPrepare("UPDATE expenses SET title = \\$2, amount = \\$3, note = \\$4, tags = \\$5 WHERE id = \\$1").WillReturnError(&pq.Error{})

		handler := Handler{DB: db}
		handler.UpdateExpensesHandler(res.Context)
		ee := util.Error{}
		res.Decode(&ee)

		assert.Equal(t, "1", res.Context.Param("id"))
		assert.Equal(t, http.StatusInternalServerError, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotNil(t, ee.Message)
	})
	t.Run("should return 500 (InternalServerError) when cannot execute update", func(t *testing.T) {
		res, db, mock := arrange("")
		mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow("1", "test-title", "123", "test-note", pq.Array([]string{"test-tags"}))
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = ?").
			ExpectQuery().
			WithArgs("1").
			WillReturnRows(mockRows)
		mock.ExpectPrepare("UPDATE expenses SET title = \\$2, amount = \\$3, note = \\$4, tags = \\$5 WHERE id = \\$1").
			ExpectExec().
			WillReturnError(&pq.Error{})

		handler := Handler{DB: db}
		handler.UpdateExpensesHandler(res.Context)
		ee := util.Error{}
		res.Decode(&ee)

		assert.Equal(t, "1", res.Context.Param("id"))
		assert.Equal(t, http.StatusInternalServerError, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotNil(t, ee.Message)
	})
}

func arrange(body string) (*util.Response, *sql.DB, sqlmock.Sqlmock) {
	e := Expense{
		ID:     1,
		Title:  "apple smoothie",
		Amount: 89,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}
	var res *util.Response
	if body == "" {
		b, _ := json.Marshal(e)
		res = util.RequestE(http.MethodPut, "/expenses/1", strings.NewReader(string(b)))
	} else {
		res = util.RequestE(http.MethodPut, "/expenses/1", strings.NewReader(body))
	}
	res.Context.SetPath("/expenses/:id")
	res.Context.SetParamNames("id")
	res.Context.SetParamValues("1")
	db, mock, _ := sqlmock.New()
	return res, db, mock
}
