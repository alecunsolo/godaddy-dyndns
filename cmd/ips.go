package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
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

type godaddyPayload struct {
	IP         string `json:"data"`
	Name       string `json:"name"`
	TTL        int    `json:"ttl"`
	RecordType string `json:"type"`
}

func currentIP() (string, error) {
	var apiURL = fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/%s", viper.GetString("domain"), viper.GetString("hostname"))
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", viper.Get("api-key"), viper.GetString("secret-key")))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	godaddyData := make([]godaddyPayload, 1)
	err = json.Unmarshal(body, &godaddyData)
	if err != nil {
		return "", err
	}
	return godaddyData[0].IP, nil
}
