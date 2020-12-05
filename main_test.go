package main

import (
	"testing"
	"log"
	"net/http"
	"net/url"
	"fmt"
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

	item := todoItem{"Call Mom"}
	requestAsJsonByteSlice, err := json.Marshal(item)

	requestAsIOWriter := bytes.NewBuffer(requestAsJsonByteSlice)
	request, err :=  http.NewRequest(
		"GET",
		"http://localhost:8080/api/findItem",
		requestAsIOWriter)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {log.Fatal(err)}
	if response.StatusCode == http.StatusAccepted {
		fmt.Println("Item was found")
	} else {
		fmt.Println("Item was not found")
	}
}

func TestAddAndDeleteItem(t *testing.T) {
	client := http.Client{}

	item := todoItem{"new test item"}
	requestBodyAsJsonByteSlice, err := json.Marshal(item)

	requestBodyAsIOWriter := bytes.NewBuffer(requestBodyAsJsonByteSlice)
	request, err :=  http.NewRequest(
		"GET",
		"http://localhost:8080/api/findItem",
		requestBodyAsIOWriter)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {log.Fatal(err)}
	if response.StatusCode == http.StatusAccepted {
		log.Fatal("The task 'new test item' was already present in the database, " +
						"obviating the test to add it")
	}

	_, err = http.PostForm("http://localhost:8080/api/newItem", url.Values{"Title": {"new test item"}})
	if err != nil {log.Fatal(err)}

	response, err = client.Do(request)
	if err != nil {log.Fatal(err)}
	if response.StatusCode != http.StatusAccepted {
		log.Fatal("Add item test was unsuccessful")
	}

	deleteRequestBodyAsIOWriter := bytes.NewBuffer(requestBodyAsJsonByteSlice)

	deleteRequest, err := http.NewRequest(
		"DELETE",
		"http://localhost:8080/api/deleteItem",
		deleteRequestBodyAsIOWriter)
	deleteRequest.Header.Set("Content-Type", "application/json")
	_, err = client.Do(deleteRequest)
	if err != nil {log.Fatal(err)}

	response, err = client.Do(request)
	if err != nil {log.Fatal(err)}
	if response.StatusCode == http.StatusAccepted {
		log.Fatal("Delete item test was unsuccessful")
	}
}

func TestAddItem(t *testing.T) {
	if testing.Short() {t.Skip()}
	_, err := http.PostForm("http://localhost:8080/api/newItem", url.Values{"Title": {"new test item"}})
	if err != nil {log.Fatal(err)}
}

func TestDeleteItem(t *testing.T) {
	if testing.Short() {t.Skip()}
	client := http.Client{}

	item := todoItem{"new test item"}
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