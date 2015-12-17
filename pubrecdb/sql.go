package pubrecdb

func createSql() string {
	// Returns the SQL command that is used to create the pubrecord.db
	// We figure out where that file is by using GOPATH

	sql := `-- DB Schema -- Version 0.2.0

CREATE TABLE blocks (
    hash        TEXT NOT NULL, 
    prevhash    TEXT UNIQUE NOT NULL, -- Unique constraint prevents forks.
    height      INT  UNIQUE NOT NULL, -- The number of blocks between this one and the genesis block.
    timestamp   INT,        -- The timestamp stored as an epoch time
    -- Extra fields added to reproduce hash of block
    version     INT,
    merkleroot  TEXT,
    difficulty  INT,        -- uint32
    nonce       INT,        -- uint32

    PRIMARY KEY(hash) -- enforces
    FOREIGN KEY (prevhash) REFERENCES blocks(hash)
);

CREATE TABLE bulletins (
    txid        TEXT NOT NULL, 
    block       TEXT NOT NULL,
    author      TEXT NOT NULL,  -- From the address of the first OutPoint used.
    message     TEXT NOT NULL,  -- UTF-8, must have some content.
    timestamp   INT,            -- Seconds since Jan 1, 1970
    latitude    INT,            -- Fixed point decimal. Divided by 1,000,000 to produce position
    longitude   INT,            -- See above
    height      INT,            -- Part of coords


    PRIMARY KEY(txid), 
    FOREIGN KEY(block) REFERENCES blocks(hash) ON DELETE CASCADE
);

CREATE TABLE endorsements (
    txid        TEXT NOT NULL, -- the enclosing transactions SHA hash
    block       TEXT NOT NULL, -- the containing block hash
    bid         TEXT NOT NULL, -- the endorsed bulletins SHA hash
    timestamp   INT NOT NULL,  -- Unix time
    author      TEXT NOT NULL, -- formatted as a bitcoin address.

    PRIMARY KEY(txid)
    FOREIGN KEY(block) REFERENCES blocks(hash) ON DELETE CASCADE
);

CREATE TABLE tags (
    txid   TEXT,
    value  TEXT,

    FOREIGN KEY(txid) REFERENCES bulletins(txid) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_tags ON tags (value);
CREATE INDEX IF NOT EXISTS idx_height ON blocks (height);
CREATE INDEX IF NOT EXISTS idx_timestamp ON blocks (timestamp);
`

	return sql
}
