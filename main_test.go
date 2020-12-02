package main

import (
	"testing"
	"log"
	"net/http"
	"net/url"
	//"net/http/httptest"
	//"encoding/json"
	//"bytes"
)

func TestAddItem(t *testing.T) {
	response, err := http.PostForm("http://localhost:8080/api/newItem", url.Values{"Title": {"new test item"}})
	if err != nil {log.Fatal(err)}
	_ = response
}