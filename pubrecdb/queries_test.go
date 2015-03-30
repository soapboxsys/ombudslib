package pubrecdb

import (
	"testing"
	"time"
)

func TestJsonBlock(t *testing.T) {

	db, _ := SetupTestDB()

	h := "00000000777213b4fd7c5d5a71b9b52608356c4194203b1b63d1bb0e6141d17d"
	jsonBlkResp, err := db.GetJsonBlock(h)

	if err != nil {
		t.Fatal(err)
	}

	respH := jsonBlkResp.Head.Hash
	if respH != h {
		t.Fatalf("Hashes don't match [%s] and returned: [%s]", h, respH)
	}

	// Check to see if empty block is reachable

	h = "00000000efaee711979fe42e667188e50b1096e4d9cfcbc9a82101336189c2ca"
	jsonBlkResp, err = db.GetJsonBlock(h)
	if err != nil {
		t.Fatal(err)
	}
	if len(jsonBlkResp.Bulletins) > 0 {
		t.Fatalf("This is an empty block")
	}
}

func TestJsonAuthor(t *testing.T) {
	db, _ := SetupTestDB()

	author := "miUDcP8obUKPhqkrBrQz57sbSg2Mz1kZXH"

	jsonResp, err := db.GetJsonAuthor(author)

	if err != nil {
		t.Fatal(err)
	}

	blkTs := int64(1414017952)
	if jsonResp.Author.NumBltns != 2 || jsonResp.Author.FirstBlkTs != blkTs {
		t.Fatalf("Wrong values:\n [%s]\n", jsonResp)
	}
}

func TestWholeBoard(t *testing.T) {
	db, _ := SetupTestDB()

	board := "ahimsa-dev"

	wholeBoard, err := db.GetWholeBoard(board)

	if err != nil {
		t.Fatal(err)
	}

	if wholeBoard.Summary.NumBltns != 4 {
		t.Fatalf("Wrong values:\n [%s]\n", wholeBoard)
	}

	expLA := int64(1414193281)
	if wholeBoard.Summary.LastActive != expLA {
		t.Fatalf(
			"Wrong last active time in:\n[%s]\nWanted an LA of: %d\n\tGot: %d",
			wholeBoard.Summary,
			expLA,
			wholeBoard.Summary.LastActive,
		)
	}

	expCA := int64(1414017952)
	if wholeBoard.Summary.CreatedAt != expCA {
		t.Fatalf(
			"Wrong created at time.\nWanted a CA of: %d\n\tGot %d",
			expCA,
			wholeBoard.Summary.LastActive,
		)
	}

}

func TestAllBoards(t *testing.T) {
	db, _ := SetupTestDB()

	allBoards, err := db.GetAllBoards()

	if err != nil {
		t.Fatal(err)
	}

	if len(allBoards) != 4 {
		t.Fatalf("Wrong number of boards returned:\n [%s]\n", allBoards)
	}

}

func TestBlockDay(t *testing.T) {
	db, _ := SetupTestDB()

	target := time.Date(2014, time.November, 1, 0, 0, 0, 0, time.UTC)

	blks, err := db.GetBlocksByDay(target)
	if err != nil {
		t.Fatal(err)
	}

	if len(blks) != 4 {
		t.Fatalf("Wrong number of blocks for this day:\n%s\n", blks)
	}
}

func TestLatestDB(t *testing.T) {
	db, _ := SetupTestDB()

	lastBlk, lastBltn, err := db.LatestBlkAndBltn()
	if err != nil {
		t.Fatal(err)
	}

	if lastBlk != 1415862580 || lastBltn != 1415854832 {
		t.Fatalf(
			"Wrong latest things returned!\nBltn: %d\nBlk: %d",
			lastBlk,
			lastBltn,
		)
	}
}
