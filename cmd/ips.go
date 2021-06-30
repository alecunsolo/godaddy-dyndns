package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type ip struct {
	IP string `json:"query"`
}

func retrieveExternalIP() (string, error) {
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
	var apiURL = fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/%s", domain, hostname)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", apiKey, keySecret))
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

type godaddyUpdatePayload struct {
	IP  string `json:"data"`
	TTL int    `json:"ttl"`
}

func updateIP(dnsIP, extIP string) error {
	if dnsIP == extIP {
		log.Println("External IP and current DNS record are equal. Nothing to do")
		return nil
	}
	log.Printf("Current external IP: %s. Current DNS record: %s", extIP, dnsIP)
	var apiURL = fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/%s", domain, hostname)
	payload := []godaddyUpdatePayload{
		{
			IP:  extIP,
			TTL: 600,
		},
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", apiKey, keySecret))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Something went wrong
	if resp.StatusCode != 200 {
		log.Printf("response Status: %s", resp.Status)
		log.Printf("response Headers: %s", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body))
	}

	return nil
}
