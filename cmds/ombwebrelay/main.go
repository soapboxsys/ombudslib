package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/soapboxsys/ombudslib/jsonapi"
	"github.com/soapboxsys/ombudslib/ombutil"
	"github.com/soapboxsys/ombudslib/pubrecdb"
)

var (
	host       = flag.String("host", "localhost:1055", "The ip and port for the server to listen on")
	pubrecpath = flag.String("pubrecpath", "", "The path to the static files to serve")
	verbose    = flag.Bool("verbose", false, "Logs the output of every request")
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()

	nodedir := ombutil.AppDataDir("ombnode", false)
	dbpath := filepath.Join(nodedir, "data", "mainnet", "pubrecord.db")

	if *pubrecpath != "" {
		dbpath = *pubrecpath
	}
	log.Printf("Opening pubrec: %s\n", dbpath)
	db, err := pubrecdb.LoadDB(dbpath)
	if err != nil {
		log.Fatal(err)
	}

	prefix := "/api/"
	router := jsonapi.Router(prefix, db)

	who := jsonapi.ApiFacts{
		Operator:      "LCD Sound Systems",
		Location:      "Anchorage, AK",
		Administrator: "AV Mike",
		ContactInst:   "Knock three times and speak Friend",
	}
	jsonapi.AddApiFacts(who, prefix, router)

	log.Printf("Webserver listening at %s.\n", *host)

	if *verbose {
		logger := Log(router)
		log.Fatal(http.ListenAndServe(*host, logger))
	} else {
		log.Fatal(http.ListenAndServe(*host, router))
	}
}
