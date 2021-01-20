package main

import (
	"os"

	"github.com/aiuzu42/SukiBot/app"
	"github.com/aiuzu42/SukiBot/config"

	log "github.com/sirupsen/logrus"
)

func main() {
	setupLog()
	err := config.InitConfig()
	if err != nil {
		log.Fatal("[main]Error loading configuration: " + err.Error())
	}
	app.StartApp()
}

func setupLog() {
	file, err := os.OpenFile("sukiLog.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
}
