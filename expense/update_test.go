package expense

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
		res, db, mock := arrange(string(b))

		mock.ExpectPrepare("UPDATE expenses SET title = \\$2, amount = \\$3, note = \\$4, tags = \\$5 WHERE id = \\$1 RETURNING id, title, amount, note, tags").
			ExpectQuery().
			WithArgs(e.ID, e.Title, e.Amount, e.Note, pq.Array(e.Tags)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
				AddRow("1", e.Title, fmt.Sprint(e.Amount), e.Note, pq.Array(e.Tags)))
		handler := Handler{DB: db}

		handler.UpdateExpensesHandler(res.Context)
		ee := Expense{}
		res.Decode(&ee)

		assert.Equal(t, "1", res.Context.Param("id"))
		assert.Equal(t, http.StatusOK, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Equal(t, e, ee)
	})
	t.Run("should return 400 (BadRequest) when request id invalid", func(t *testing.T) {
		res := util.RequestE(http.MethodPut, "/expenses/invalid", nil)
		res.Context.SetPath("/expenses/:id")
		res.Context.SetParamNames("id")
		res.Context.SetParamValues("invalid")
		db, mock, _ := sqlmock.New()

		handler := Handler{DB: db}
		handler.UpdateExpensesHandler(res.Context)
		ee := util.Error{}
		res.Decode(&ee)

		assert.Equal(t, "invalid", res.Context.Param("id"))
		assert.Equal(t, http.StatusBadRequest, res.Recorder.Code)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotNil(t, ee.Message)
	})
	t.Run("should return 400 (BadRequest) when request body is invalid", func(t *testing.T) {
		res, db, mock := arrange("invalid body")

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
		mock.ExpectPrepare("UPDATE expenses SET title = \\$2, amount = \\$3, note = \\$4, tags = \\$5 WHERE id = \\$1 RETURNING id, title, amount, note, tags").
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
	t.Run("should return 500 (InternalServerError) when cannot prepare update", func(t *testing.T) {
		res, db, mock := arrange("")
		mock.ExpectPrepare("UPDATE expenses SET title = \\$2, amount = \\$3, note = \\$4, tags = \\$5 WHERE id = \\$1 RETURNING id, title, amount, note, tags").
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
	t.Run("should return 500 (InternalServerError) when cannot execute update", func(t *testing.T) {
		res, db, mock := arrange("")
		mock.ExpectPrepare("UPDATE expenses SET title = \\$2, amount = \\$3, note = \\$4, tags = \\$5 WHERE id = \\$1 RETURNING id, title, amount, note, tags").
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
