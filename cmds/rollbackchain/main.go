// This script connects to a Bitcoin peer requests a block based on the
// "target" block hash and attempts to write it to a file.
package main

import (
	"flag"
	"log"
	"path"

	"github.com/btcsuite/btcd/database"
	_ "github.com/btcsuite/btcd/database/ldb"
	"github.com/btcsuite/btcd/wire"
	"github.com/soapboxsys/ombudslib/pubrecdb"
)

var defaultDir = "/home/ubuntu/.ombnode/data"
var dataDir = flag.String("datadir", defaultDir, "The nodes data directory")
var testNet = flag.Bool("testnet", false, "The testnet switch")
var targetStr = flag.String("target", "", "The hash of the first block to be dropped")

func main() {
	flag.Parse()

	target, err := wire.NewShaHashFromStr(*targetStr)
	if err != nil {
		log.Fatal(err)
	}

	dataPath := path.Join(*dataDir, "mainnet")
	if *testNet {
		dataPath = path.Join(*dataDir, "testnet")
	}

	// Load pubrec
	p := path.Join(dataPath, "pubrecord.db")
	precdb, err := pubrecdb.LoadDB(p)
	if err != nil {
		log.Fatal("Loading prec: ", err)
	}

	// Load Block chain
	dbPath := path.Join(dataPath, "blocks_leveldb")
	db, err := database.OpenDB("leveldb", dbPath)
	defer db.Close()
	if err != nil {
		log.Fatal("Loading leveldb: ", err)
	}

	printDBStates(db, precdb)

	t, err := db.FetchBlockHeaderBySha(target)
	if err != nil {
		log.Fatal("LevelDB fetch: ", err)
	}

	log.Println("Dropping leveldb blocks....")

	// Disconnect all blocks after Sha
	if err := db.DropAfterBlockBySha(&t.PrevBlock); err != nil {
		log.Print("Leveldb Drop Threw: ", err)
	}

	log.Println("Dropping pubrecord blocks....")
	// Delete from pubrec up to shahash.
	if err = precdb.DropAfterBlockBySha(&t.PrevBlock); err != nil {
		log.Print("Pubrec Drop Threw: ", err)
	}

	log.Println("Deleted all blocks after target sha in both DBs")
	printDBStates(db, precdb)
}

func printDBStates(db database.Db, prec *pubrecdb.PublicRecord) {
	newest, h, err := db.NewestSha()
	if err != nil {
		log.Fatal("Leveldb newest: ", err)
	}

	log.Printf("Leveldb newest sha at height %d is: %s", h, newest.String())
	tip, err := prec.GetBlockTip()
	if err != nil {
		log.Fatal("Pubrec get: ", err)
	}
	log.Printf("Pubrecdb newest sha at height %d is: %s", tip.Head.Height, tip.Head.Hash)
}
