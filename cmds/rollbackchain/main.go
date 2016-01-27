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

var node = flag.String("node", "", "The node to connect to")
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
	log.Println(p)
	precdb, err := pubrecdb.LoadDB(p)
	if err != nil {
		log.Fatal("Loading prec: ", err)
	}

	// Load Block chain
	dbPath := path.Join(dataPath, "blocks_leveldb")
	db, err := database.OpenDB("leveldb", dbPath)
	if err != nil {
		log.Fatal("Loading leveldb: ", err)
	}

	t, err := db.FetchBlockHeaderBySha(target)
	if err != nil {
		log.Fatal("LevelDB fetch: ", err)
	}

	// Disconnect all blocks after Sha
	if err := db.DropAfterBlockBySha(&t.PrevBlock); err != nil {
		log.Print("Leveldb Drop Threw: ", err)
	}

	// Delete from pubrec up to shahash.
	if err = precdb.DropAfterBlockBySha(&t.PrevBlock); err != nil {
		log.Print("Pubrec Drop Trew: ", err)
	}

	log.Println("Deleted all blocks after target sha")
	newest, h, err := db.NewestSha()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Leveldb newest sha at height %d is: %s", h, newest.String())

	tip, err := precdb.GetBlockTip()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("PubrecordDb newest sha at height %d is %s", tip.Head.Height, tip.Head.Hash)
}
