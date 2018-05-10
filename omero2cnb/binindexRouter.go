package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nimezhu/data"
)

type BinindexRouter struct {
	index  map[string]*data.BinIndexMap
	mem    map[string]*AnnotationMapValue
	dbname string
	omero  string //omero url
}

func (db *BinindexRouter) ServeTo(router *mux.Router) {
	router.HandleFunc("/genomes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		genomes := []string{}
		for k, _ := range db.index {
			genomes = append(genomes, k)
		}
		t, _ := json.Marshal(&genomes)
		w.Write(t)
	})
	router.HandleFunc("/{genome}/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.WriteString(w, "[\"omero\"]")
	})
	router.HandleFunc("/{genome}/ls", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		params := mux.Vars(r)
		genome := params["genome"]
		io.WriteString(w, "[{\"dbname\":\"omero\",\"format\":\"binindex\",\"genome\":\""+genome+"\",\"uri\":\"null\",\"omero\":\""+db.omero+"\"}]")
	})
	router.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {

	})
	router.HandleFunc("/{genome}/"+db.dbname+"/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//id := params["id"]
		io.WriteString(w, "[{\"id\":\"omero\",\"format\":\"binindex\",\"omero\":\""+db.omero+"\"}]")
	})
	router.HandleFunc("/{genome}/"+db.dbname+"/omero/get/{chr}:{start}-{end}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		params := mux.Vars(r)
		//id := params["id"]
		genome := params["genome"]
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
		if idx, ok := db.index[genome]; ok {
			vals, err := idx.QueryRegion(chrom, s, e)
			if err != nil {
				io.WriteString(w, "{'error':'not found'}")
				return
			}
			for v := range vals {
				p, _ := db.mem[v.(*BedURI).ParentID()]
				fmt.Println(p.AnnotationID, p.ParentType)
				io.WriteString(w, fmt.Sprintf("%s\t%d\t%d\t%s\t%s\t%d\t%s\n", chrom, v.Start(), v.End(), v.Id(), p.ParentType, p.ParentID, v.(*BedURI).Color()))
			}
		}
	})
}
