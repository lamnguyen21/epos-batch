package main

import (
	"fmt"
	"testing"
)

func TestItemProcessing(t *testing.T) {
	itemDAO := ItemDAO{}
	err := itemDAO.init("../data_26_oct_18.xlsx")
	if err != nil {
		t.Error(err)
	} else {
		itemDAO.processItems(func(items []Item) {
			fmt.Println(items)
		})
	}
}
