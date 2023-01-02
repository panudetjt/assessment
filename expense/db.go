package expense

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/panudetjt/assessment/util"
)

func InitDB() *util.DB {
	url := os.Getenv("DATABASE_URL")

	var err error
	rdb, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	db := &util.DB{DB: rdb}
	return db
}
