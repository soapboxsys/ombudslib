package jsonapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/gorilla/mux"
	"github.com/soapboxsys/ombudslib/ombjson"
	"github.com/soapboxsys/ombudslib/ombutil"
	"github.com/soapboxsys/ombudslib/pubrecdb"
)

var (
	processStart time.Time = time.Now()
)

func writeJson(w http.ResponseWriter, m interface{}) {

	bytes, err := json.Marshal(m)
	if err != nil {
		http.Error(w, "Failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func BulletinHandler(db *pubrecdb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		txidStr, _ := mux.Vars(request)["txid"]
		txid, err := wire.NewShaHashFromStr(txidStr)
		if err != nil {
			http.Error(w, "That is not a sha2 hash", 404)
			return
		}

		bltn, err := db.GetBulletin(txid)
		if err == sql.ErrNoRows {
			http.Error(w, "Bulletin does not exist", 404)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, bltn)
	}
}

func EndorsementHandler(db *pubrecdb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		txidStr, _ := mux.Vars(request)["txid"]
		txid, err := wire.NewShaHashFromStr(txidStr)
		if err != nil {
			http.Error(w, "That is not a sha2 hash", 404)
			return
		}

		bltn, err := db.GetEndorsement(txid)
		if err == sql.ErrNoRows {
			http.Error(w, "Endorsement does not exist", 404)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, bltn)
	}
}

func BlockHandler(db *pubrecdb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		hashStr, _ := mux.Vars(request)["hash"]
		hash, err := wire.NewShaHashFromStr(hashStr)
		if err != nil {
			http.Error(w, "That is not a sha2 hash", 404)
			return
		}

		bltn, err := db.GetBlock(hash)
		if err == sql.ErrNoRows {
			http.Error(w, "Block is not in record", 404)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, bltn)
	}
}

// Handles serving a bulletin board.
func TagHandler(db *pubrecdb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		tagstr, _ := mux.Vars(request)["tag"]
		tag := ombutil.Tag("#" + tagstr)

		board, err := db.GetTag(tag)
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), 405)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, board)
	}
}

func NewHandler(db *pubrecdb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		page, err := db.GetLatestPage()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, page)
	}
}

func StatusHandler(db *pubrecdb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		blk, err := db.GetBlockTip()
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		stat := &ombjson.Status{
			BlockTip: blk.Head,
		}
		writeJson(w, stat)
	}
}

// returns the http handler initialized with the api's routes. The prefix should
// start and end with slashes. For example /api/ is a good prefix.
func Handler(prefix string, db *pubrecdb.PublicRecord) http.Handler {

	r := mux.NewRouter()
	sha2re := "([a-f]|[A-F]|[0-9]){64}"
	//addrgex := "([a-z]|[A-Z]|[0-9]){30,35}"
	// Since the tag's path could be percent encoded we give it 3x wiggle room
	// since a single byte in percent encoding is %EE.
	tagre := ".{1,90}"

	// A single day follows this format: DD-MM-YY
	//dayre := `[0-9]{1,2}-[0-9]{1,2}-[0-9]{4}`

	p := prefix
	// Item handlers
	r.HandleFunc(p+fmt.Sprintf("bltn/{txid:%s}", sha2re), BulletinHandler(db))
	r.HandleFunc(p+fmt.Sprintf("endo/{txid:%s}", sha2re), EndorsementHandler(db))
	r.HandleFunc(p+fmt.Sprintf("block/{hash:%s}", sha2re), BlockHandler(db))

	// Aggregate handlers
	r.HandleFunc(p+fmt.Sprintf("tag/{tag:%s}", tagre), TagHandler(db))
	r.HandleFunc(p+"new", NewHandler(db))

	// Meta handlers
	r.HandleFunc(p+"status", StatusHandler(db))

	return r
}
