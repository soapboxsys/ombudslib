package pubrecdb

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"

	"github.com/mxk/go-sqlite/sqlite3"
)

// The overarching struct that contains everything needed for a connection to a
// sqlite db containing the public record.
type PublicRecord struct {
	conn *sql.DB
	// Read only connection for filtering
	roConn *sqlite3.Conn

	// Precompiled SQL selects
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

	// Precompiled inserts
	insertBlockHeadStmt   *sql.Stmt
	insertBulletinStmt    *sql.Stmt
	insertTagStmt         *sql.Stmt
	insertEndorsementStmt *sql.Stmt

	// Precompiled deletes
	deleteBlockStmt *sql.Stmt

	// Utility queries
	blockIsTipStmt *sql.Stmt
}

// Creates a DB at the desired path or drops an existing one and recreates a
// new empty one at the path. The bitcoin network is needed because the genesis
// block must be inserted first for the DB to initialized properly.
func InitDB(path string, params *chaincfg.Params) (*PublicRecord, error) {
	path = filepath.Clean(path)
	// Check if the file exists and remove it if it does.
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			return nil, err
		}
	}

	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Execute the Create Table sql
	_, err = conn.Exec(createSql())
	if err != nil {
		return nil, err
	}
	conn.Close()

	db, err := createPubRec(path)
	if err != nil {
		return nil, err
	}

	if err := prepareInserts(db); err != nil {
		return nil, err
	}

	// Insert the net's Genesis block
	genesisBlk := btcutil.NewBlock(params.GenesisBlock)
	genesisBlk.SetHeight(0)

	if ok, err := db.InsertBlockHead(genesisBlk); !ok {
		return nil, err
	}

	return prepareDB(db)
}

// Loads a sqlite db, checks if its reachabale and prepares all the queries.
func LoadDB(path string) (*PublicRecord, error) {
	db, err := createPubRec(path)
	if err != nil {
		return nil, err
	}
	return prepareDB(db)
}

func createPubRec(path string) (*PublicRecord, error) {
	path = filepath.Clean(path)
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	roPath := "file://" + path + "?mode=ro"
	roConn, err := sqlite3.Open(roPath)
	if err != nil {
		return nil, err
	}

	db := &PublicRecord{
		conn:   conn,
		roConn: roConn,
	}

	return db, nil
}

// prepareDB takes a pubrecord and initializes all of the precompiled
// statements and executes any connection specific code
func prepareDB(db *PublicRecord) (*PublicRecord, error) {

	if err := execPragma(db); err != nil {
		return nil, fmt.Errorf("Pragma defs failed: %s", err)
	}

	/*if err := prepareQueries(db); err != nil {
		return nil, fmt.Errorf("Preparing queries failed: %v", err)
	}*/

	if err := prepareInserts(db); err != nil {
		return nil, fmt.Errorf("Preparing inserts failed: %v", err)
	}

	if err := prepareDeletes(db); err != nil {
		return nil, fmt.Errorf("Preparing deletes failed: %v", err)
	}

	return db, nil
}

// execPragma executes directives that are needed for the write side of the SQL
// conn to enforce high quality (and secure!) sql statements.
func execPragma(db *PublicRecord) error {
	// The following pragmas define the operation of the sqlite3 conn. This
	// does important things: it enforces foreign key constraints, ...
	pragmas := `
	PRAGMA foreign_keys=ON;
	`

	if _, err := db.conn.Exec(pragmas); err != nil {
		return err
	}
	return nil
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

	return nil
}
