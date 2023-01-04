package expense

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/panudetjt/assessment/util"
	"github.com/stretchr/testify/assert"
)

func TestGetExpenseByIdHandler(t *testing.T) {
	t.Run("should return 200 (OK) when request is valid", func(t *testing.T) {
		res := util.RequestE(http.MethodGet, "/expenses/1", nil)
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
		handler := Handler{DB: db}

		handler.GetExpenseByIdHandler(res.Context)

		assert.Equal(t, "1", res.Context.Param("id"))
		assert.Equal(t, http.StatusOK, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return 404 (NotFound) when no item", func(t *testing.T) {
		res := util.RequestE(http.MethodGet, "/expenses/0", nil)
		res.Context.SetPath("/expenses/:id")
		res.Context.SetParamNames("id")
		res.Context.SetParamValues("0")
		db, mock, _ := sqlmock.New()
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = ?").
			ExpectQuery().
			WithArgs("0").
			WillReturnError(sql.ErrNoRows)
		handler := Handler{DB: db}

		handler.GetExpenseByIdHandler(res.Context)

		assert.Equal(t, http.StatusNotFound, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return 500 (InternalServerError) when cannot prepare", func(t *testing.T) {
		res := util.RequestE(http.MethodGet, "/expenses/0", nil)
		res.Context.SetPath("/expenses/:id")
		res.Context.SetParamNames("id")
		res.Context.SetParamValues("0")

		db, mock, _ := sqlmock.New()
		handler := Handler{DB: db}

		handler.GetExpenseByIdHandler(res.Context)
		var e util.Error
		res.Decode(&e)

		assert.Equal(t, http.StatusInternalServerError, res.Recorder.Code)
		assert.NotEmpty(t, e.Message)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return 500 (InternalServerError) when cannot scan", func(t *testing.T) {
		res := util.RequestE(http.MethodGet, "/expenses/1", nil)
		res.Context.SetPath("/expenses/:id")
		res.Context.SetParamNames("id")
		res.Context.SetParamValues("0")

		db, mock, _ := sqlmock.New()
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = ?").
			ExpectQuery().
			WithArgs("0")
		handler := Handler{DB: db}

		handler.GetExpenseByIdHandler(res.Context)
		var e util.Error
		res.Decode(&e)

		assert.Equal(t, http.StatusInternalServerError, res.Recorder.Code)
		assert.NotEmpty(t, e.Message)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
