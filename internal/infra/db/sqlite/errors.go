package sqlite

import (
	"errors"

	"modernc.org/sqlite"
	"modernc.org/sqlite/lib"
)

// IsUniqueViolation reports whether err is SQLite rejecting a
// duplicate value in a UNIQUE column.
func IsUniqueViolation(err error) bool {
	var sqliteErr *sqlite.Error
	return errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
}
