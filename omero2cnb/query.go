package main

import (
	"database/sql"
	"fmt"
	"strconv"
)

func getParentIDType(annotationID int, db *sql.DB) (int, string, bool) {
	var t = []string{"image", "well", "project", "dataset"}
	for _, v := range t {
		if idx, ok := _getParentID(annotationID, db, v); ok {
			fmt.Println(idx, v)
			return idx, v, true
		}
	}
	return -1, "", false
}
func getImageParentID(id int, db *sql.DB) (int, bool) {
	return _getParentID(id, db, "image")
}

func _getParentID(annotationID int, db *sql.DB, t string) (int, bool) {
	sql := "SELECT parent FROM " + t + "annotationlink WHERE child=" + strconv.Itoa(annotationID)
	rows, err := db.Query(sql)
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
