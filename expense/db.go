package expense

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func InitDB(driver string) (*sql.DB, error) {
	url := os.Getenv("DATABASE_URL")

	db, err := sql.Open(driver, url)
	if err != nil {
		return nil, fmt.Errorf("connect to database error: %s", err)
	}

	return db, nil
}
