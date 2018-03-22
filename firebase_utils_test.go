package main

import (
	"testing"
)

func TestInitFireBase(t *testing.T) {
	app, err := InitFirebase()

	if err != nil {
		t.Error(err)
	} else {
		t.Log(app)
	}
}

func TestDatabase(t *testing.T) {
	app, _ := InitFirebase()
	client, err := Database(app)

	if err != nil {
		t.Error(err)
	} else {
		t.Log(client)
	}
}

func TestGetItem(t *testing.T) {
	app, _ := InitFirebase()
	client, err := Database(app)

	if err != nil {
		t.Error(err)
	} else {
		items, err := GetItem(client)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(items)
		}
	}
}

func TestProcessItems(t *testing.T) {
	app, _ := InitFirebase()
	client, err := Database(app)

	if err != nil {
		t.Error(err)
	} else {
		processItems(client)
	}
}
