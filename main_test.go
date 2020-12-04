package main

import (
	"testing"
	"log"
	"net/http"
	"net/url"
	//"io/ioutil"
	//"net/http/httptest"
	"encoding/json"
	"bytes"
)

func TestAddItem(t *testing.T) {
	_, err := http.PostForm("http://localhost:8080/api/newItem", url.Values{"Title": {"new test item"}})
	if err != nil {log.Fatal(err)}
}

func TestDeleteItem(t *testing.T) {
	client := http.Client{}

	item := todoItem{"Go to store"}
	jsonReq, err := json.Marshal(item)
	request, err :=  http.NewRequest(
		"DELETE",
		"http://localhost:8080/api/deleteItem",
		bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {log.Fatal(err)}

	if response.Status != "202 Accepted" {
		log.Fatal("Response was not 202 Accepted; it was ", response.Status)
	}
}