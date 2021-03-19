package helpers

import "log"

func LogError(message string, err error) {
	if err != nil {
		log.Println(message, err.Error())
	}
}
