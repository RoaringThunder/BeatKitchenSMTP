package main

import (
	"fmt"
	"log"
	"salamander-smtp/app"
	"salamander-smtp/database"
	"salamander-smtp/logging"
)

func init() {
	err := database.InitializeThunderDome()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Welllcomeeee to the Thunderrrrr Dome!")

	err = logging.InitLogging()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Hello World")
	err := app.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
