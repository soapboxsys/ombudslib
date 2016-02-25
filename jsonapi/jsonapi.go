package jsonapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
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

func NewStatsHandler(db *pubrecdb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		now := time.Now()
		// Look back one day
		yesterday := now.Add(-(time.Hour * 24))
		stats, err := db.GetStatistics(yesterday, now)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJson(w, stats)
	}
}

func BestTagsHandler(db *pubrecdb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		tags, err := db.GetBestTags()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, tags)
	}
}

func AuthorHandler(db *pubrecdb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		addrStr, _ := mux.Vars(request)["addr"]
		// Try our best to decode the passed AddrStr. If we can't parse it.
		// Drop it.
		var author btcutil.Address
		var err error
		author, err = btcutil.DecodeAddress(addrStr, &chaincfg.MainNetParams)
		if err != nil {
			author, err = btcutil.DecodeAddress(addrStr, &chaincfg.TestNet3Params)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
		resp, err := db.GetAuthor(author)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, resp)
	}
}

func NearbyLocHandler(db *pubrecdb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		latStr, _ := mux.Vars(request)["lat"]
		lonStr, _ := mux.Vars(request)["lon"]
		rStr, _ := mux.Vars(request)["r"]

		var err error
		var lat, lon, r float64

		if lat, err = strconv.ParseFloat(latStr, 64); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if lon, err = strconv.ParseFloat(lonStr, 64); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if r, err = strconv.ParseFloat(rStr, 64); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		bltns, err := db.GetNearbyBltns(lat, lon, r)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, bltns)
	}
}

func MostEndoHandler(db *pubrecdb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		bltns, err := db.GetMostEndorsedBltns(10)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, bltns)
	}
}

func StatusHandler(db *pubrecdb.PublicRecord, start time.Time) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		blk, err := db.GetBlockTip()
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		uptime := time.Now().Sub(start)
		ts := time.Now()

		stat := &ombjson.Status{
			BlockTip:   blk.Head,
			Uptime:     int64(uptime.Seconds()),
			UptimeH:    uptime.String(),
			Timestamp:  ts.Unix(),
			TimestampH: ts.String(),
		}
		writeJson(w, stat)
	}
}

type ApiFacts struct {
	UserAgent     string `json:"user-agent",omitempty`
	Operator      string `json:"operator",omitempty`
	Location      string `json:"location",omitempty`
	Administrator string `json:"admin",omitempty`
	ContactInst   string `json:"HowToReach",omitempty`
}

func AddApiFacts(who ApiFacts, prefix string, router *mux.Router) {
	f := func(w http.ResponseWriter, request *http.Request) {
		who.UserAgent = "ombudslib/jsonapi"
		writeJson(w, who)
	}
	router.HandleFunc(prefix+"whoami", f)
}

func Handler(prefix string, db *pubrecdb.PublicRecord) http.Handler {
	return Router(prefix, db)
}

// returns the http handler initialized with the api's routes. The prefix should
// start and end with slashes. For example /api/ is a good prefix.
func Router(prefix string, db *pubrecdb.PublicRecord) *mux.Router {

	r := mux.NewRouter()
	sha2re := "([a-f]|[A-F]|[0-9]){64}"
	addrgex := "([a-z]|[A-Z]|[0-9]){30,35}"
	// Since the tag's path could be percent encoded we give it 3x wiggle room
	// since a single byte in percent encoding is %EE.
	tagre := ".{1,90}"

	// A single day follows this format: DD-MM-YY
	//dayre := `[0-9]{1,2}-[0-9]{1,2}-[0-9]{4}`

	// Pulls floats out of urls
	l_re := `[-+]?(\d*[.])?\d+`
	loc_suffix := fmt.Sprintf("loc/{lat:%s},{lon:%s},{r:%s}", l_re, l_re, l_re)

	p := prefix
	// Item handlers
	r.HandleFunc(p+fmt.Sprintf("bltn/{txid:%s}", sha2re), BulletinHandler(db))
	r.HandleFunc(p+fmt.Sprintf("endo/{txid:%s}", sha2re), EndorsementHandler(db))
	r.HandleFunc(p+fmt.Sprintf("block/{hash:%s}", sha2re), BlockHandler(db))
	r.HandleFunc(p+fmt.Sprintf("author/{addr:%s}", addrgex), AuthorHandler(db))
	r.HandleFunc(p+loc_suffix, NearbyLocHandler(db))

	// Paginated handlers
	r.HandleFunc(p+fmt.Sprintf("tag/{tag:%s}", tagre), TagHandler(db))
	r.HandleFunc(p+"new", NewHandler(db))

	// Aggregate handlers
	r.HandleFunc(p+"pop-tags", BestTagsHandler(db))
	r.HandleFunc(p+"most-endo", MostEndoHandler(db))

	// Meta handlers
	r.HandleFunc(p+"status", StatusHandler(db, time.Now()))
	r.HandleFunc(p+"new/statistics", NewStatsHandler(db))

	return r
}
