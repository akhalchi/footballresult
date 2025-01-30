package main

import (
	"footballresult/config"
	"footballresult/get"
	"footballresult/send"
	"footballresult/storage"
	"log"
	"sync"
)

func main() {
	config.LoadConfig()

	db, err := storage.InitDB()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("database is connected")
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		get.Events(db)
	}()

	go func() {
		defer wg.Done()
		send.Telegram(db)
	}()

	wg.Wait()

}
