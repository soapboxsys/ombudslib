package pubrecdb

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	_ "github.com/mattn/go-sqlite3"
	"github.com/soapboxsys/ombudslib/ombwire/peg"
)

var defaultMaxQueryLimit = 10000

// The overarching struct that contains everything needed for a connection to a
// sqlite db containing the public record.
type PublicRecord struct {
	conn *sql.DB

	// Max Number of Records returned by db
	maxQueryLimit int

	// Precompiled SQL selects
	selectBltn        *sql.Stmt
	selectTag         *sql.Stmt
	selectEndo        *sql.Stmt
	selectBltnsHeight *sql.Stmt
	findHeight        *sql.Stmt
	selectBlock       *sql.Stmt
	selectBlockTip    *sql.Stmt
	selectEndosByBid  *sql.Stmt
	selectBestTags    *sql.Stmt
	selectAuthorBtlns *sql.Stmt
	selectAuthorEndos *sql.Stmt

	// Line-O-PROGRESS
	selectBlockHead   *sql.Stmt
	selectBlockBltns  *sql.Stmt
	selectAuthor      *sql.Stmt
	selectAuthorBltns *sql.Stmt
	selectBlacklist   *sql.Stmt
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

	// Prepare the DB to do one insert (for the pegBlk)
	if err := prepareInserts(db); err != nil {
		return nil, err
	}

	if params.Net == wire.MainNet {
		// Insert the pegged starting block
		pegBlk := peg.GetStartBlock()
		if err, ok := db.InsertBlockHead(pegBlk); !ok || err != nil {
			return nil, err
		}
	} else if params.Net == wire.TestNet3 {
		pegBlk := peg.GetTestStartBlock()
		if err, ok := db.InsertBlockHead(pegBlk); !ok || err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("No peg for non-default Bitcoin Net")
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

	db := &PublicRecord{
		conn: conn,
	}

	return db, nil
}

// prepareDB takes a pubrecord and initializes all of the precompiled
// statements and executes any connection specific code
func prepareDB(db *PublicRecord) (*PublicRecord, error) {

	db.maxQueryLimit = defaultMaxQueryLimit

	if err := execPragma(db); err != nil {
		return nil, fmt.Errorf("Pragma defs failed: %s", err)
	}

	if err := prepareQueries(db); err != nil {
		return nil, fmt.Errorf("Preparing queries failed: %v", err)
	}

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
