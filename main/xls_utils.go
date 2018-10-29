package main

import (
	"github.com/tealeg/xlsx"
)

type ItemDAO struct {
	sheet *xlsx.Sheet
}

func (itemDAO *ItemDAO) init(xlsFile string) error {
	file, err := xlsx.OpenFile(xlsFile)
	if err != nil {
		return err
	} else {
		itemDAO.sheet = file.Sheets[0]
		return nil
	}
}

func (itemDAO *ItemDAO) processItems(callback func([]Item)) {
	buffer := make([]Item, 10)
	for _, row := range itemDAO.sheet.Rows {
		if len(buffer) < 10 {
			buffer = append(buffer, Item{
				Code:           row.Cells[1].String(),
				VmName:         row.Cells[6].String(),
				RetailPriceVAT: row.Cells[10].String(),
				ExtraPrice:     row.Cells[11].String(),
			})
		} else {
			callback(buffer)
			buffer = buffer[len(buffer):]
		}
	}
}

func (itemDAO *ItemDAO) processOriginalItems(callback func([]Item)) {
	buffer := make([]Item, 10)
	for _, row := range itemDAO.sheet.Rows {
		if len(buffer) < 10 {
			buffer = append(buffer, Item{
				Code:           row.Cells[1].String(),
				VmName:         row.Cells[6].String(),
				RetailPriceVAT: row.Cells[14].String(),
			})
		} else {
			callback(buffer)
			buffer = buffer[len(buffer):]
		}
	}
}
