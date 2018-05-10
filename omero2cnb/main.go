package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/nimezhu/data"
)

var dbmem map[string]*AnnotationMapValue
var dbindex map[string]*data.BinIndexMap
var db *sql.DB

const AppName string = "Omerome Browser"
const VERSION string = "0.0.1"
const coordsKey string = "regions"

var app App = App{AppName, VERSION}

func dbindexInsert(d *AnnotationMapValue) {
	if d.Name == coordsKey {
		prefix := strconv.Itoa(d.Index)
		k := getKey(d)
		if arr, ok := parseRegions(d.Value, prefix, k); ok {
			for _, v := range arr {
				if _, ok1 := dbindex[v.Genome()]; !ok1 {
					dbindex[v.Genome()] = data.NewBinIndexMap()
				}
				dbindex[v.Genome()].Insert(v)
			}
		}
	}
}
func dbindexDelete(d *AnnotationMapValue) {
	if d.Name == coordsKey {
		prefix := strconv.Itoa(d.Index)
		k := getKey(d)
		if arr, ok := parseRegions(d.Value, prefix, k); ok {
			for _, v := range arr {
				dbindex[v.Genome()].Delete(v)
			}
		}
	}
}
func getKey(d *AnnotationMapValue) string {
	return fmt.Sprintf("%d:%d", d.AnnotationID, d.Index)
}
func insertEntry(d *AnnotationMapValue) error {
	k := getKey(d)
	dbmem[k] = d
	dbindexInsert(d)
	return nil
}
func deleteEntry(d *AnnotationMapValue) error {
	delete(dbmem, getKey(d))
	dbindexDelete(d)
	return nil
}
func updateEntry(d *AnnotationMapValue) error {
	k := getKey(d)
	previous := dbmem[k]
	if d.Name == coordsKey {
		dbindexDelete(previous)
		dbindexInsert(d)
	}
	delete(dbmem, k)
	dbmem[k] = d
	return nil
}
func addParentTypeID(a *RawAnnotationMapValue) *AnnotationMapValue {
	if idx, t, ok := getParentIDType(a.AnnotationID, db); ok {
		d := AnnotationMapValue{a.AnnotationID, a.Name, a.Value, a.Index, t, idx}
		return &d
	} else {
		d := AnnotationMapValue{a.AnnotationID, a.Name, a.Value, a.Index, "unknown", -1}
		return &d
	}
}
func waitForNotification(l *pq.Listener) {
	for {
		select {
		case n := <-l.Notify:
			fmt.Println("Received data from channel [", n.Channel, "] :")
			// Prepare notification payload for pretty print
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
			if err != nil {
				fmt.Println("Error processing JSON: ", err)
				return
			}
			//fmt.Println(string(prettyJSON.Bytes()))
			//TODO update dbmem
			var action Action
			json.Unmarshal(prettyJSON.Bytes(), &action)
			fmt.Println(action.Action, action.Data.Name)
			processedData := addParentTypeID(&action.Data)
			if action.Action == "INSERT" {
				insertEntry(processedData)
			} else if action.Action == "DELETE" {
				deleteEntry(processedData)
			} else if action.Action == "UPDATE" {
				updateEntry(processedData)

			}
			return
		case <-time.After(90 * time.Second):
			fmt.Println("Received no events for 90 seconds, checking connection")
			go func() {
				l.Ping()
			}()
			return
		}
	}
}

/* BinIndexing Server, Connect to db
 *   server $DBNAME　$USER　$PASSWD
 */
func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	host := os.Args[1]
	dbname := os.Args[2]
	user := os.Args[3]
	passwd := os.Args[4]
	omero := os.Args[5] //omero server
	conninfo := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", host, dbname, user, passwd)
	port := 3721
	var err error
	db, err = sql.Open("postgres", conninfo)
	dbmem = map[string]*AnnotationMapValue{}
	dbindex = map[string]*data.BinIndexMap{}
	defer db.Close()

	if err != nil {
		panic(err)
	}
	//TODO Processing Available Data
	rows, err := db.Query("SELECT * FROM annotation_mapvalue")
	checkErr(err)
	annos := []*RawAnnotationMapValue{}
	for rows.Next() {
		var annotationID int
		var name string
		var value string
		var index int
		err = rows.Scan(&annotationID, &name, &value, &index)
		checkErr(err)
		fmt.Printf("%3v | %8v | %6v | %6v\n", annotationID, name, value, index)
		a := RawAnnotationMapValue{annotationID, name, value, index}
		annos = append(annos, &a)
	}
	for _, a := range annos {
		insertEntry(addParentTypeID(a))
	}

	//manager

	//TODO Serve HTTP dbmem
	router := mux.NewRouter()
	//add manager
	//manager := Manager{dbmem, ""}
	//manager.ServeTo(router)
	binManager := BinindexRouter{dbindex, dbmem, "omero", omero}
	binManager.ServeTo(router)

	router.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		a, err2 := json.Marshal(app)
		if err2 == nil {
			w.Write(a)
		} else {
			w.Write([]byte("{'error':'not found'}"))
		}
	})

	go http.ListenAndServe(":"+strconv.Itoa(port), router)

	//TODO Process Updating Data
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(conninfo, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("events")
	if err != nil {
		panic(err)
	}

	fmt.Println("Start monitoring PostgreSQL...")
	for {
		waitForNotification(listener)
	}
}
