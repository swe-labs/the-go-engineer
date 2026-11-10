// Package pgdriver registers the pgx/v5 stdlib driver under the "pgx"
// database/sql driver name.
//
// The exercise directories _starter/ and _solution/ are main packages, but
// directories beginning with an underscore are excluded from the go tool's
// package walk — including go mod tidy. Importing the driver through this
// visible package keeps the github.com/jackc/pgx/v5 requirement resolvable
// so `go mod tidy` stays idempotent and the project builds in CI.
package pgdriver

import (
	_ "github.com/jackc/pgx/v5/stdlib"
)
