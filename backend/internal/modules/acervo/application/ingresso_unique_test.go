package application

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestUniqueViolation(t *testing.T) {
	if uniqueViolation(nil) || uniqueViolation(errors.New("x")) {
		t.Fatal("erro comum não é 23505")
	}
	if !uniqueViolation(&pgconn.PgError{Code: "23505"}) {
		t.Fatal("unicidade Postgres")
	}
}
