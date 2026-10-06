package dbtest

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"workshop/internal/infra/db/sqlc"
)

// FailingQuerier returns a db.Querier whose every method returns err.
//
// It is the generated *db.Queries over a database that can never connect,
// so every query sqlc generates fails with err without a method being
// written for it here.
func FailingQuerier(err error) db.Querier {
	return db.New(sql.OpenDB(failingConnector{err: err}))
}

// failingConnector is a driver.Connector whose every connection fails
// with err, which database/sql hands back unwrapped from each query.
type failingConnector struct {
	err error
}

// Connect implements driver.Connector.
func (c failingConnector) Connect(context.Context) (driver.Conn, error) {
	return nil, c.err
}

// Driver implements driver.Connector.
func (c failingConnector) Driver() driver.Driver {
	return failingDriver(c)
}

// failingDriver is the driver.Driver behind failingConnector.
type failingDriver failingConnector

// Open implements driver.Driver.
func (d failingDriver) Open(string) (driver.Conn, error) {
	return nil, d.err
}
