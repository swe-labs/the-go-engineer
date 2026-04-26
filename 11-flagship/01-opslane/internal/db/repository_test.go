package db

import (
	"fmt"
	"testing"

	"github.com/lib/pq"
)

func TestIsForeignKeyViolation(t *testing.T) {
	t.Parallel()

	err := fmt.Errorf("wrapped: %w", &pq.Error{Code: postgresForeignKeyViolation})

	if !isForeignKeyViolation(err) {
		t.Fatal("expected wrapped PostgreSQL foreign key violation to be detected")
	}
}
