package main

import (
	"testing"
	"log"
	"net/http"
	"net/url"
	//"fmt"
	//"io/ioutil"
	//"net/http/httptest"
	"encoding/json"
	"bytes"
 	//"os"
	//"database/sql"
	//"io/ioutil"
)

func TestFindItem(t *testing.T) {
	client := http.Client{}

	item := todoItem{"Call Dad"}
	jsonReq, err := json.Marshal(item)

	bytesReq := bytes.NewBuffer(jsonReq)
	request, err :=  http.NewRequest(
		"GET",
		"http://localhost:8080/api/findItem",
		bytesReq)
	request.Header.Set("Content-Type", "application/json")
	_, err = client.Do(request)
	if err != nil {log.Fatal(err)}
}

func TestAddAndDeleteItem(t *testing.T) {
	
}

func TestAddItem(t *testing.T) {
	if testing.Short() {t.Skip()}
	_, err := http.PostForm("http://localhost:8080/api/newItem", url.Values{"Title": {"new test item"}})
	if err != nil {log.Fatal(err)}
}

func TestDeleteItem(t *testing.T) {
	if testing.Short() {t.Skip()}
	client := http.Client{}

	item := todoItem{"Call Mom"}
	jsonReq, err := json.Marshal(item)
	//fmt.Printf( "%s\n", jsonReq)

	bytesReq := bytes.NewBuffer(jsonReq)
	request, err :=  http.NewRequest(
		"DELETE",
		"http://localhost:8080/api/deleteItem",
		bytesReq)
	request.Header.Set("Content-Type", "application/json")
	_, err = client.Do(request)
	if err != nil {log.Fatal(err)}
}