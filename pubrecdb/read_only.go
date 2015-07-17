package pubrecdb

import (
	"errors"

	"github.com/soapboxsys/ombudslib/ombjson"
)

// ExecuteROExpr uses the db's read-only connection to execute a select query
// that only returns rows of bulletins. This command is run the read-only connection
// to the sqlite db and the query is checked by the sqlite3_stmt_readonly function.
func (db *PublicRecord) ExecuteROExpr(expr string) ([]*ombjson.JsonBltn, error) {
	// Check Expr

	empty := []*ombjson.JsonBltn{}

	stmt, err := db.roConn.Prepare(expr)
	if err != nil {
		return empty, err
	}

	if !stmt.ReadOnly() {
		return empty, errors.New("Statement attempts to modify the DB!")
	}

	err = stmt.Query()
	if err != nil {
		return empty, err
	}

	/*bltns, err := getRelevantBltns(rows)
	if err != nil {
		return empty, err
	}*/

	return empty, err
}
