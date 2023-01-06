package expense

import (
	"database/sql"
	"errors"
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

func TestAllExpenseHandler(t *testing.T) {
	t.Run("should return 200 (OK) when request is valid", func(t *testing.T) {
		p := prepare()
		res := p.res
		mock := p.mock
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses").
			ExpectQuery().
			WillReturnRows(p.mockRows)
		handler := Handler{DB: p.db}

		handler.GetAllExpenseHandler(res.Context)
		var es []Expense
		res.Decode(&es)

		assert.Equal(t, http.StatusOK, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Equal(t, 2, len(es))
		assert.Equal(t, p.expenses, es)
	})

	t.Run("should return 200 (OK) even no row in database", func(t *testing.T) {
		p := prepare()
		res := p.res
		mock := p.mock
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses").
			ExpectQuery().
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}))
		handler := Handler{DB: p.db}

		handler.GetAllExpenseHandler(res.Context)
		var es []Expense
		res.Decode(&es)

		assert.Equal(t, http.StatusOK, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Equal(t, 0, len(es))
	})

	t.Run("should return 500 (InternalServerError) when cannot prepare SELECT", func(t *testing.T) {
		p := prepare()
		res := p.res
		mock := p.mock
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses").
			WillReturnError(errors.New("error"))
		handler := Handler{DB: p.db}

		handler.GetAllExpenseHandler(res.Context)
		var e util.Error
		res.Decode(&e)

		assert.Equal(t, http.StatusInternalServerError, res.Recorder.Code)
		assert.NotEmpty(t, e.Message)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return 500 (InternalServerError) when cannot execute SELECT", func(t *testing.T) {
		p := prepare()
		res := p.res
		mock := p.mock
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses").
			ExpectQuery().
			WillReturnError(&pq.Error{})
		handler := Handler{DB: p.db}

		handler.GetAllExpenseHandler(res.Context)
		var e util.Error
		res.Decode(&e)

		assert.Equal(t, http.StatusInternalServerError, res.Recorder.Code)
		assert.NotEmpty(t, e.Message)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("should return 500 (InternalServerError) when cannot SELECT scan", func(t *testing.T) {
		p := prepare()
		res := p.res
		mock := p.mock
		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses").
			ExpectQuery().
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		handler := Handler{DB: p.db}

		handler.GetAllExpenseHandler(res.Context)
		var e util.Error
		res.Decode(&e)

		assert.Equal(t, http.StatusInternalServerError, res.Recorder.Code)
		assert.NotEmpty(t, e.Message)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

type prepared struct {
	res      *util.Response
	mockRows *sqlmock.Rows
	db       *sql.DB
	mock     sqlmock.Sqlmock
	expenses []Expense
}

func prepare() prepared {
	es := []Expense{
		{
			ID:     1,
			Title:  "test-title-1",
			Amount: 123,
			Note:   "test-note-2",
			Tags:   []string{"test-tags-3"},
		},
		{
			ID:     2,
			Title:  "test-title-1",
			Amount: 456,
			Note:   "test-note-2",
			Tags:   []string{"test-tags-3"},
		},
	}
	res := util.RequestE(http.MethodGet, "/expenses", nil)
	res.Context.SetPath("/expenses")
	mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"})
	for _, e := range es {
		mockRows.AddRow(e.ID, e.Title, e.Amount, e.Note, pq.Array(e.Tags))
	}
	db, mock, _ := sqlmock.New()
	return prepared{res, mockRows, db, mock, es}
}
