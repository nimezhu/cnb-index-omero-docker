package main

import (
	"database/sql"
	"strconv"
)

func getParentIDType(annotationID int, db *sql.DB) (int, string, bool) {
	var t = []string{"image", "well", "project", "dataset"}
	for _, v := range t {
		if idx, ok := _getParentID(annotationID, db, v); ok {
			return idx, v, true
		}
	}
	return -1, "", false
}
func getImageParentID(id int, db *sql.DB) (int, bool) {
	return _getParentID(id, db, "image")
}

func _getParentID(annotationID int, db *sql.DB, t string) (int, bool) {
	rows, err := db.Query("SELECT parent FROM " + t + "annotationlink where child=" + strconv.Itoa(annotationID))
	parentID := -1
	if err != nil {
		return parentID, false
	}
	sign := false
	for rows.Next() {
		err = rows.Scan(&parentID)
		if err == nil {
			sign = true
		}
	}
	return parentID, sign
}
