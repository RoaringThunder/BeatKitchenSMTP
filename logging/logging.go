package logging

import (
	"fmt"
	"log"
	"os"
)

func InitLogging() error {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	fmt.Println("Initialized logging")
	return err
}

func LogError(err string, resp uint) {
	log.Println(fmt.Sprintf("Error: %s", err))
	fmt.Println(fmt.Sprintf("Error: %s", err))
}

func Log(err string) {
	log.Println(err)
	fmt.Println(err)
}
