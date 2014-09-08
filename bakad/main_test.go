package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func init() {
	setup(1 * time.Minute)
}

func TestAnnounce(t *testing.T) {
	defer reset()

	values := url.Values{}
	values.Add("url", "http://test.com:1234")
	body := strings.NewReader(values.Encode())

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/announce", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	http.DefaultServeMux.ServeHTTP(resp, req)
	if resp.Code != 200 {
		t.Error("should return HTTP OK", resp.Code, "≠", 200)
	}

	var result []string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"http://test.com:1234"}
	if !reflect.DeepEqual(result, expected) {
		t.Error("unexpected response", result, "≠", expected)
	}
}

func reset() {
	peerlist = newPeers(1 * time.Minute)
}
