package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	item := todoItem{Title: itemTitle}

	err := database.findItem(itemTitle)
	if err == nil {t.Error("item ", itemTitle, "was already in database," +
		"obviating the test to add it to the database")}

	err = database.addItem(item)
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

	err := findItemViaAPICall(item)
	if err == nil {
		t.Errorf("Task with title %v was already present in the database, "+
			"obviating the test to add it", item.Title)
	}

	err = insertItemViaAPICall(itemTitle, item)
	if err != nil {t.Error(err)}

	err = findItemViaAPICall(item)
	if err != nil {
		t.Error("Item", item.Title, "was not found after it was added")
	}

	err = deleteItemViaAPICall(item)
	if err != nil {t.Error(err)}

	err = findItemViaAPICall(item)
	if err == nil {
		t.Error("Item", item.Title, "was still in the database " +
			"when the test should have deleted it")
	}
}

func deleteItemViaAPICall(item todoItem) error {
	client := http.Client{}
	err, deleteRequest := setupDeleteRequest(item)
	_, err = client.Do(deleteRequest)
	return err
}


func insertItemViaAPICall(itemTitle string, item todoItem) error {
	_, err := http.PostForm(
		"http://localhost:8080/api/newItem",
		url.Values{"Title": {itemTitle}})
	if err != nil {return err}
	return err
}

func findItemViaAPICall(item todoItem) error {
	client := http.Client{}
	request, _ := setupFindRequest(item)
	response, _ := client.Do(request)
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("item %v was not found", item.Title)
	}
	return nil
}

func setupFindRequest(item todoItem) (*http.Request, error) {
	requestBodyAsJsonByteSlice, _ := json.Marshal(item)
	requestBodyAsIOWriter := bytes.NewBuffer(requestBodyAsJsonByteSlice)
	request, err := http.NewRequest(
		"GET",
		"http://localhost:8080/api/findItem",
		requestBodyAsIOWriter)
	request.Header.Set("Content-Type", "application/json")
	return request, err
}

func setupDeleteRequest(item todoItem) (error, *http.Request) {
	err, deleteRequestBodyAsIOWriter := convertToJSONInIOWriter(item)
	deleteRequest, err := http.NewRequest(
		"DELETE",
		"http://localhost:8080/api/deleteItem",
		deleteRequestBodyAsIOWriter)
	deleteRequest.Header.Set("Content-Type", "application/json")
	return err, deleteRequest
}

func convertToJSONInIOWriter(item todoItem) (error, *bytes.Buffer) {
	requestBodyAsJsonByteSlice, err := json.Marshal(item)
	deleteRequestBodyAsIOWriter := bytes.NewBuffer(requestBodyAsJsonByteSlice)
	return err, deleteRequestBodyAsIOWriter
}

func printTodoList(t *testing.T) {
	todoList := database.getTodoList()
	fmt.Println("===========Todo Items ============")
	for _, item := range todoList {
		fmt.Println(item.Title)
	}
	fmt.Println("===========End list===============")
}

