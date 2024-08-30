package main

import (
	"fmt"
	"footballresult/get"
	"footballresult/storage"
	"log"
)

func main() {
	db, err := storage.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db)

	get.Matches()

}
