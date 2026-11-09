package main

import (
	"errors"
	"fmt"
)

type DBError struct {
	Code    string
	Message string
	Err     error
}

func (e *DBError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("DB %s: %s (cause: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("DB %s: %s", e.Code, e.Message)
}

func (e *DBError) Unwrap() error {
	return e.Err
}

func (e *DBError) Is(target error) bool {
	t, ok := target.(*DBError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

var ErrNotFound = &DBError{Code: "NOT_FOUND", Message: "not found"}

func queryDB(id string) error {
	switch {
	case id == "":
		return fmt.Errorf("queryDB: %w", &DBError{Code: "NOT_FOUND", Message: "user not found"})
	case id == "dup":
		return fmt.Errorf("queryDB: %w", &DBError{
			Code:    "DUPLICATE_KEY",
			Message: "duplicate",
			Err:     errors.New("unique constraint violation"),
		})
	default:
		return nil
	}
}

func main() {
	for _, id := range []string{"", "dup", "valid"} {
		err := queryDB(id)
		var dbErr *DBError
		switch {
		case errors.As(err, &dbErr):
			fmt.Printf("queryDB(%q): code=%s msg=%q", id, dbErr.Code, dbErr.Message)
			if dbErr.Err != nil {
				fmt.Printf(" cause=%v", dbErr.Err)
			}
			fmt.Println()
		case err == nil:
			fmt.Printf("queryDB(%q): success\n", id)
		default:
			fmt.Printf("queryDB(%q): unexpected: %v\n", id, err)
		}

		if errors.Is(err, ErrNotFound) {
			fmt.Printf("  -> matches ErrNotFound\n")
		}
	}
}
