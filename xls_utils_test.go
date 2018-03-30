package main

import (
	"fmt"
	"testing"
)

func TestItemProcessing(t *testing.T) {
	itemDAO := ItemDAO{}
	err := itemDAO.init("data.xlsx")
	if err != nil {
		t.Error(err)
	} else {
		itemDAO.processItems(func(items []Item) {
			fmt.Println(items)
		})
	}
}
