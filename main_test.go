package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestHandlers(t *testing.T) {
	database = newDatabaseConnection()
	t.Run("Print", printTodoList)
	t.Run("Add and delete", AddAndDeleteOneItem)
}

func TestHandlersViaAPICalls(t *testing.T) {
	t.Run("Add and delete", AddAndDeleteAnItemViaAPICalls)
}

func AddAndDeleteOneItem(t *testing.T) {
	itemTitle := "dummy test item"

	err := database.findItem(itemTitle)
	if err == nil {t.Error("item ", itemTitle, "was already in database," +
		"obviating the test to add it to the database")}

	err = database.addItem(itemTitle)
	if err != nil {t.Error(err)}

	err = database.findItem(itemTitle)
	if err != nil {t.Error(err)}

	err = database.deleteItem(itemTitle)
	if err != nil {t.Error(err)}

	err = database.findItem(itemTitle)
	if err == nil {t.Error("item with title", itemTitle, "was still in the database " +
		"when it should have been deleted")}
}

func AddAndDeleteAnItemViaAPICalls(t *testing.T) {
	const itemTitle = "new test item"
	item := todoItem{itemTitle}

	err := findItem(item)
	if err == nil {
		t.Errorf("Task with title %v was already present in the database, "+
			"obviating the test to add it", item.Title)
	}

	err = addItem(itemTitle, item)
	if err != nil {t.Error(err)}

	err = findItem(item)
	if err != nil {
		t.Error("Item", item.Title, "was not found after it was added")
	}

	err = deleteItem(item)
	if err != nil {t.Error(err)}

	err = findItem(item)
	if err == nil {
		t.Error("Item", item.Title, "was still in the database " +
			"when the test should have deleted it")
	}
}

func deleteItem(item todoItem) error {
	client := http.Client{}
	requestBodyAsJsonByteSlice, err := json.Marshal(item)
	deleteRequestBodyAsIOWriter := bytes.NewBuffer(requestBodyAsJsonByteSlice)
	deleteRequest, err := http.NewRequest(
		"DELETE",
		"http://localhost:8080/api/deleteItemHandler",
		deleteRequestBodyAsIOWriter)
	deleteRequest.Header.Set("Content-Type", "application/json")
	_, err = client.Do(deleteRequest)
	return err
}

func addItem(itemTitle string, item todoItem) error {
	_, err := http.PostForm(
		"http://localhost:8080/api/newItem",
		url.Values{"Title": {itemTitle}})
	if err != nil {return err}
	return err
}

func findItem(item todoItem) error {
	client := http.Client{}
	requestBodyAsJsonByteSlice, _ := json.Marshal(item)
	requestBodyAsIOWriter := bytes.NewBuffer(requestBodyAsJsonByteSlice)
	request, _ := http.NewRequest(
		"GET",
		"http://localhost:8080/api/findItemHandler",
		requestBodyAsIOWriter)
	request.Header.Set("Content-Type", "application/json")
	response, _ := client.Do(request)
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("item %v was not found", item.Title)
	}
	return nil
}

func NewFindRequest(item todoItem) {

}

func TestPingDatabase(t *testing.T) {
	err := database.ping()
	if err != nil {t.Error(err)}
}

func printTodoList(t *testing.T) {
	todoList := database.getTodoList()
	fmt.Println("===========Todo Items ============")
	fmt.Println(strings.Join(todoList, "\n"))
	fmt.Println("===========End list===============")
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
		"http://localhost:8080/api/deleteItemHandler",
		bytesReq)
	request.Header.Set("Content-Type", "application/json")
	_, err = client.Do(request)
	if err != nil {log.Fatal(err)}
}

