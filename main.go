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

	// Используем sync.WaitGroup для синхронизации горутин
	var wg sync.WaitGroup

	// Увеличиваем счетчик на 2, так как будут запущены 2 горутины
	wg.Add(2)

	// Запускаем функцию get.Events параллельно
	go func() {
		defer wg.Done() // Уменьшаем счетчик при завершении горутины
		get.Events(db)
	}()

	// Запускаем функцию send.Telegram параллельно
	go func() {
		defer wg.Done() // Уменьшаем счетчик при завершении горутины
		send.Telegram(db)
	}()

	// Ожидаем завершения обеих горутин
	wg.Wait()

	log.Println("All parallel tasks completed.")
}
