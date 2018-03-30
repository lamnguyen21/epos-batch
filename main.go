package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

func main() {
	var path string
	var data string

	flag.StringVar(&path, "file", "", "Specify a path to your serviceAccount.json")
	flag.StringVar(&data, "data", "", "Specify a path to your data excel file")
	flag.Parse()

	if strings.Compare(path, "") == 0 || strings.Compare(data, "") == 0 {
		panic("Missing service account file or data file")
	}

	file, err := os.Create("batch.log")
	if err != nil {
		panic(err)
	}

	logger := log.New(file, "", log.LstdFlags|log.Lshortfile)
	fbStore := FireBaseStore{}

	if err := fbStore.init(path); err != nil {
		panic(err)
	} else {
		itemDAO := ItemDAO{}
		err := itemDAO.init(data)
		if err != nil {
			panic(err)
		} else {
			itemDAO.processItems(func(items []Item) {
				for _, item := range items {
					logger.Printf("Update item with code %v\n", item.Code)
					fbStore.updateItemPrice(item)
				}
			})
		}
	}
}
