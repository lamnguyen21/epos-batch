package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type LogonResponse struct {
	LoginName  string
	SessionKey string
}

type ItemResponse struct {
	ObjectName          string
	ObjectNumber        string
	UnitPrice           float64
	DiscountedUnitPrice float64
}

func logon() (*LogonResponse, error) {
	res, err := http.PostForm("https://resume.little-demo.tk/pos/api/Session/LogOn",
		url.Values{"Username": {"katy.le"}, "Password": {"1"}})

	if err != nil {
		return nil, err
	}

	b, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var logonReponse LogonResponse
	err = json.Unmarshal(b, &logonReponse)
	if err != nil {
		return nil, err
	}

	return &logonReponse, nil
}

func getItem(itemCode string, logonResponse *LogonResponse) (*ItemResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://resume.little-demo.tk/pos/api/itemUnitPriceByObjectNumber/"+itemCode, nil)
	req.SetBasicAuth(logonResponse.LoginName, logonResponse.SessionKey)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	b, _ := ioutil.ReadAll(res.Body)
	var item ItemResponse
	err = json.Unmarshal(b, &item)

	if err != nil {
		return nil, err
	}

	return &item, err
}
