package pubrecdb

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/soapboxsys/ombudslib/protocol/ombproto"
)

// The overarching struct that contains everything needed for a connection to a
// sqlite db containing the public record
type PublicRecord struct {
	conn *sql.DB

	// Precompiled SQL for ombprotorest
	selectTxid        *sql.Stmt
	selectBlockHead   *sql.Stmt
	selectBlockBltns  *sql.Stmt
	selectAuthor      *sql.Stmt
	selectAuthorBltns *sql.Stmt
	selectBlacklist   *sql.Stmt
	selectBoardSum    *sql.Stmt
	selectBoardBltns  *sql.Stmt
	selectAllBoards   *sql.Stmt
	selectRecentConf  *sql.Stmt
	selectUnconfirmed *sql.Stmt
	selectBlksByDay   *sql.Stmt
	selectDBStatus    *sql.Stmt
	selectAllAuthors  *sql.Stmt

	// Precompiled Inserts
	insertBlock *sql.Stmt
}

// Loads a sqlite db, checks if its reachabale and prepares all the queries.
func LoadDB(path string) (*PublicRecord, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	db := &PublicRecord{
		conn: conn,
	}

	err = prepareQueries(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Creates a DB at the desired path or drops an existing one and recreates a
// new empty one at the path.
func InitDB(path string) (*PublicRecord, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Get the database schema for the public record.
	create, err := ombproto.GetCreateSql()
	if err != nil {
		return nil, err
	}

	dropcmd := `
	DROP TABLE IF EXISTS blocks;
	DROP TABLE IF EXISTS bulletins;
	DROP TABLE IF EXISTS blacklist;
	`

	// DROP db if it exists and recreate it.
	_, err = conn.Exec(dropcmd + create)
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	return LoadDB(path)
}

// Prepares all of the selects for maximal speediness note that all of the queries
// must be initialized here or nil pointers will get thrown at runtime.
func prepareQueries(db *PublicRecord) error {

	var err error
	db.selectTxid, err = db.conn.Prepare(selectTxidSql)
	if err != nil {
		return err
	}

	db.selectBlockHead, err = db.conn.Prepare(selectBlockHeadSql)
	if err != nil {
		return err
	}

	db.selectBlockBltns, err = db.conn.Prepare(selectBlockBltnsSql)
	if err != nil {
		return err
	}

	db.selectAuthor, err = db.conn.Prepare(selectAuthorSql)
	if err != nil {
		return err
	}

	db.selectAuthorBltns, err = db.conn.Prepare(selectAuthorBltnsSql)
	if err != nil {
		return err
	}

	db.selectBlacklist, err = db.conn.Prepare(selectBlacklistSql)
	if err != nil {
		return err
	}

	db.selectBoardSum, err = db.conn.Prepare(selectBoardSumSql)
	if err != nil {
		return err
	}

	db.selectBoardBltns, err = db.conn.Prepare(selectBoardBltnsSql)
	if err != nil {
		return err
	}

	db.selectBoardSum, err = db.conn.Prepare(selectBoardSumSql)
	if err != nil {
		return err
	}

	db.selectBoardBltns, err = db.conn.Prepare(selectBoardBltnsSql)
	if err != nil {
		return err
	}

	db.selectAllBoards, err = db.conn.Prepare(selectAllBoardsSql)
	if err != nil {
		return err
	}

	db.selectRecentConf, err = db.conn.Prepare(selectRecentConfSql)
	if err != nil {
		return err
	}

	db.selectUnconfirmed, err = db.conn.Prepare(selectUnconfirmedSql)
	if err != nil {
		return err
	}

	db.selectBlksByDay, err = db.conn.Prepare(selectBlksByDaySql)
	if err != nil {
		return err
	}

	db.selectDBStatus, err = db.conn.Prepare(selectDBStatusSql)
	if err != nil {
		return err
	}

	db.selectAllAuthors, err = db.conn.Prepare(selectAllAuthors)
	if err != nil {
		return err
	}

	// Compile inserts
	db.insertBlock, err = db.conn.Prepare(insertBlock)
	if err != nil {
		return err
	}

	return nil
}

func SetupTestDB() (*PublicRecord, error) {

	var dbpath string

	testEnvPath := os.Getenv("TEST_DB_PATH")
	if testEnvPath != "" {
		dbpath = testEnvPath
	} else {
		dbpath = os.Getenv("GOPATH") + "/src/github.com/soapboxsys/ombudslib/pubrecdb/test.db"
		dbpath = filepath.Clean(dbpath)
	}
	var err error
	db, err := LoadDB(dbpath)
	if err != nil {
		return nil, err
	}
	return db, nil
}
