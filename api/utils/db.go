package utils

import "database/sql"

func CheckDBStatus(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil

}
