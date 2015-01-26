-- DB Schema -- Version 0.1.2


CREATE TABLE blocks (
    hash        TEXT NOT NULL, 
    prevhash    TEXT, 
    height      INT,        -- The number of blocks between this one and the genesis block.
    timestamp   INT,        -- The timestamp stored as an epoch time

    PRIMARY KEY(hash)
    FOREIGN KEY(prevhash) REFERENCES blocks(hash)
);

CREATE TABLE bulletins (
    author      TEXT NOT NULL,  -- From the address of the first OutPoint used.
    txid        TEXT NOT NULL, 
    board       TEXT,           -- UTF-8
    message     TEXT NOT NULL,  -- UTF-8, must have some content.
    timestamp   INT,            -- Seconds since Jan 1, 1970
    block       TEXT,

    PRIMARY KEY(txid), 
    FOREIGN KEY(block) REFERENCES blocks(hash)
);

-- The point of the blacklist is to highlight the fact that editorial control is still possible,
-- but now the choice is given explicity to the third party.
create TABLE blacklist ( 
    txid    TEXT,
    reason  TEXT NOT NULL,

    PRIMARY KEY(txid),
    FOREIGN KEY(txid) REFERENCES bulletins(txid)
);

CREATE INDEX IF NOT EXISTS idx_height ON blocks (height);
CREATE INDEX IF NOT EXISTS idx_timestamp ON blocks (timestamp);
