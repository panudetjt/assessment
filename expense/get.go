package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/panudetjt/assessment/util"
)

func (h *Handler) GetExpenseByIdHandler(c echo.Context) error {
	id := c.Param("id")
	stmt, err := h.DB.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.Error{Message: "can't prepare query user statment:" + err.Error()})
	}
	row := stmt.QueryRow(id)
	ep := Expense{}
	err = row.Scan(&ep.ID, &ep.Title, &ep.Amount, &ep.Note, pq.Array(&ep.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, util.Error{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, ep)
	default:
		return c.JSON(http.StatusInternalServerError, util.Error{Message: "can't scan expense:" + err.Error()})
	}
}
