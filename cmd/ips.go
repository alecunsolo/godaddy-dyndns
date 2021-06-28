package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ip struct {
	IP string `json:"query"`
}

func extIP() (string, error) {
	params := url.Values{}
	params.Add("fields", "query")

	req, err := http.NewRequest("GET", "http://ip-api.com/json/", nil)
	if err != nil {
		return "", err
	}
	req.URL.Query().Add("fields", "string")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ipData := ip{}
	err = json.Unmarshal(body, &ipData)
	if err != nil {
		return "", err
	}
	return ipData.IP, nil

}
