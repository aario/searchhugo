package main

import (
    "database/sql"
    "log"
	_ "github.com/mattn/go-sqlite3"
)

const (
    DB_FILENAME = "searchhugo-thesaurus.db"
    DB_PATH_1 = "/usr/local/share"
    DB_PATH_2 = "/usr/share"
)

var (
    db *sql.DB
)

func getDbFilePath() string {
    path := getExecutablePath() + "/" + DB_FILENAME
    if (fileExists(path)) {
        return path
    }

    path = DB_PATH_1 + "/" + DB_FILENAME
    if (fileExists(path)) {
        return path
    }

    path = DB_PATH_2 + "/" + DB_FILENAME
    if (fileExists(path)) {
        return path
    }

    panic("Failed to find DB file " + DB_FILENAME)
}
    
func dbConnect() {
    dbFilePath := getDbFilePath()

    var err error
    db, err = sql.Open("sqlite3", dbFilePath)
    if err != nil {
        panic("Failed to open sqlite database from " + dbFilePath + ":\n" + err.Error())
    }
}

func dbDisconnect() {
    db.Close()
}

func dbSelect(query string) [][]string {
	var (
		result    [][]string
		container []string
		pointers  []interface{}
	)

    rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	length := len(cols)

	for rows.Next() {
		pointers = make([]interface{}, length)
		container = make([]string, length)

		for i := range pointers {
			pointers[i] = &container[i]
		}

		err = rows.Scan(pointers...)
		if err != nil {
			panic(err.Error())
		}

		result = append(result, container)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}
