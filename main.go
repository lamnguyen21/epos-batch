package main

import (
	"flag"
	"log"
	"strings"
)

func main() {
	var path string
	flag.StringVar(&path, "file", "", "Specify a path to your serviceAccount.json")
	flag.Parse()

	if strings.Compare(path, "") == 0 {
		panic("Missing service account file")
	}

	app, _ := InitFirebaseWithAccountFile(path)
	client, err := Database(app)

	if err != nil {
		log.Panicf("Unable to connect to firebase %v", err.Error())
	} else {
		log.Println("Start processing items...")
		processItems(client)
	}
}
