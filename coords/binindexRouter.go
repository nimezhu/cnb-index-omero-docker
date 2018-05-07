package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nimezhu/data"
)

type BinindexRouter struct {
	index  *data.BinIndexMap
	dbname string //sheet tab name
}

func (db *BinindexRouter) ServeTo(router *mux.Router) {
	/* TODO
	router.HandleFunc("/"+db.dbname+"/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("id\tpos\n"))
		for k, v := range db.dataMap {
			io.WriteString(w, fmt.Sprintf("%s\t%s\n", k, bedsText(v.Position)))
		}
	})
	router.HandleFunc("/"+db.dbname+"/ls", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		s, _ := json.Marshal(db.idToUri)
		w.Write(s)
	})
	*/
	router.HandleFunc("/"+db.dbname+"/get/{chr}:{start}-{end}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		params := mux.Vars(r)
		//id := params["id"]
		chrom := params["chr"]
		s, err := strconv.Atoi(params["start"])
		if err != nil {
			io.WriteString(w, "{'error':'wrong format'}")
			return
		}
		e, err := strconv.Atoi(params["end"])
		if err != nil {
			io.WriteString(w, "{'error':'wrong format'}")
			return
		}
		vals, err := db.index.QueryRegion(chrom, s, e)
		if err != nil {
			io.WriteString(w, "{'error':'not found'}")
			return
		}
		for v := range vals {
			io.WriteString(w, fmt.Sprintf("%s\t%d\t%d\t%s", chrom, v.Start(), v.End(), v.Id()))
			io.WriteString(w, "\n")
		}
	})
}
