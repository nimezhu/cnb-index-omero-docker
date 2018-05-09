package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Manager struct {
	Data   map[string]*AnnotationMapValue
	prefix string
}

func (m *Manager) ServeTo(r *mux.Router) {
	r.HandleFunc(m.prefix+"/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		keys := []string{}
		for _, v := range m.Data {
			keys = append(keys, v.Key)
		}
		jsonHic, _ := json.Marshal(keys)
		w.Write(jsonHic)
	})
}
