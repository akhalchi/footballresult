package main

import (
	"footballresult/send"
)

func main() {
	/*db, err := storage.InitDB()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(db)
	}

	get.Events(db) */
	send.Telegram()
}
