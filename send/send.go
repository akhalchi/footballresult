package send

import (
	"footballresult/send/telegram"
	"log"
)

func Telegram() {

	message := "Villarreal 1 - 5 Barça\nPrimera Division\n"

	err := telegram.SendMessageToTelegram(message)
	if err != nil {
		log.Panic(err)
	} else {
		log.Println("Сообщение успешно отправлено.")
	}
}
