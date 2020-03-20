package saarflagid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

const defaultStatusURL = "https://scoreboard.ctf.saarland/attack.json"

func GetIDsFromStatus(service, teamIP string, data []byte) ([]string, error) {
	var s struct {
		FlagIDs map[string]map[string]map[string][]string `json:"flag_ids"`
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	serviceInfo, ok := s.FlagIDs[service]
	if !ok {
		return nil, fmt.Errorf("no such service: %s", service)
	}

	teamInfo, ok := serviceInfo[teamIP]
	if !ok {
		return nil, fmt.Errorf("no team with such IP: %s", teamIP)
	}

	var flagIDs []string

	for _, v := range teamInfo {
		flagIDs = append(flagIDs, v...)
	}

	sort.Strings(flagIDs)

	return flagIDs, nil
}

func GetIDsFromURL(service, teamIP, url string) ([]string, error) {
	c := http.Client{
		Transport: http.DefaultTransport,
		Timeout:   10 * time.Second,
	}

	resp, err := c.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch status: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch status: %v", err)
	}

	return GetIDsFromStatus(service, teamIP, data)
}

func GetIDs(service, teamIP string) ([]string, error) {
	return GetIDsFromURL(service, teamIP, defaultStatusURL)
}
