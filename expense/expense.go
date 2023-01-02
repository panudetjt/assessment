package expense

import "github.com/panudetjt/assessment/util"

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount int      `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type ExpenseHandler struct {
	DB util.IDB
}
