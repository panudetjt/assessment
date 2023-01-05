package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/panudetjt/assessment/util"
)

func (h *Handler) UpdateExpensesHandler(c echo.Context) error {
	id := c.Param("id")
	stmt, err := h.DB.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.Error{Message: "can't prepare query expense statment:" + err.Error()})
	}

	row := stmt.QueryRow(id)
	e := Expense{}
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, util.Error{Message: "expense not found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.Error{Message: "can't scan expense:" + err.Error()})
	}

	err = c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.Error{Message: err.Error()})
	}

	stmt, err = h.DB.Prepare("UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.Error{Message: "can't prepare update expense statment:" + err.Error()})
	}
	_, err = stmt.Exec(e.ID, e.Title, e.Amount, e.Note, pq.Array(e.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.Error{Message: "can't execute update expense statment:" + err.Error()})
	}

	return c.JSON(http.StatusOK, e)
}
