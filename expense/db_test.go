package expense

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDB(t *testing.T) {
	os.Setenv("DATABASE_URL", "url")
	t.Run("can connect", func(t *testing.T) {
		db, err := InitDB("postgres")

		assert.Nil(t, err)
		assert.NotNil(t, db)
	})
	t.Run("cannot connect", func(t *testing.T) {
		db, err := InitDB("invalid")

		assert.NotNil(t, err)
		assert.Nil(t, db)
	})
}
