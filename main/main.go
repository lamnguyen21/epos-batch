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
	var mode uint

	flag.StringVar(&path, "file", "", "Specify a path to your serviceAccount.json")
	flag.StringVar(&data, "data", "", "Specify a path to your data excel file")
	flag.UintVar(&mode, "mode", 0, "Specify a run mode. 0 for updating item price, 1 for updating vendor price")
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
					if mode == 0 {
						logger.Printf("Update item with code %v\n", item.Code)
						fbStore.updateItemPrice(item)
					} else {
						logger.Printf("Update vendor price with code %v\n", item.Code)
						fbStore.updateVendorPrice(item)
					}
				}
			})
		}
	}
}
