package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	

	_ "github.com/mattn/go-sqlite3"
)

func openDB(dbfile string) (*sql.DB) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}

	return(db)
}

func getHandleID(db *sql.DB, handle string) ([]int){
	var handleID []int

	rows, err := db.Query("SELECT rowid, id, service FROM handle where id=?", handle )
	// fmt.Println("SELECT rowid, id, service FROM handle where id='%s'", handle)
	if err != nil {
		log.Fatal(err)
	}
	defer 	rows.Close()

	for rows.Next() {
		var rowid int
		var id string
		var service string
		err = rows.Scan(&rowid, &id, &service)
		if err != nil {
			log.Fatal(err)
		}
		handleID = append(handleID, rowid)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return handleID
}

func getMessages(db *sql.DB, handles []int, name string) {
	for _, handle := range handles {

		query := `
			select 
			text, 
			datetime(message.date/1000000000 + strftime("%s", "2001-01-01") ,"unixepoch","localtime"),
			is_from_me 
			from message 
			WHERE handle_id=? 
			ORDER by date
		`
		rows, err := db.Query(query, handle)
		// fmt.Println("SELECT rowid, id, service FROM handle where id='%s'", handle)
		if err != nil {
			log.Fatal(err)
		}
		defer 	rows.Close()

		for rows.Next() {
			var text string
			var date string
			var is_from_me int
			rows.Scan(&text, &date, &is_from_me)

			if is_from_me == 1 {
				fmt.Printf("%s,%s,%s\n", date, "Me", text)
			} else {
				fmt.Printf("%s,%s,%s\n", date, name, text)
			}
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	// fmt.Println("Running exporter")

	if len(os.Args) != 4 {
		fmt.Println("Usage: messages_export <dbfile> <phonenumber> <fromname>")
		fmt.Println("Where phonenumber is a phone number starting with + or an email address")
		fmt.Println("Fromname is just a nicely formated name for your output")
		os.Exit(1)
	}
	dbfile := os.Args[1]
	handle := os.Args[2]
	name := os.Args[3]

	var db *sql.DB
	db = openDB(dbfile)
	var handles = getHandleID(db, handle)
	getMessages(db, handles, name)

	db.Close()
}
