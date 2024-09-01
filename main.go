package main

import (
	"footballresult/get"
	"footballresult/storage"
	"log"
)

func main() {
	db, err := storage.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	get.Matches(db)

}
