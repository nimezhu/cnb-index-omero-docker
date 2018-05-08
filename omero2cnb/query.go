package main

import (
	"database/sql"
	"errors"
	"strconv"
)

/* TODO */
func getParentIdType(annotationID int, db *sql.DB) error {
	return errors.New("TODO")
}

func getImageParentID(annotationID int, db *sql.DB) (int, bool) {
	rows, err := db.Query("SELECT parent FROM imageannotationlink where child=" + strconv.Itoa(annotationID))
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
