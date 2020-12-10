package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
)

type todoItem struct {
	Title string
}

var (
	database *DatabaseConnection
)

func getTodoListHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	todoItems := database.getTodoList()
	json.NewEncoder(writer).Encode(todoItems)
}

func addItemHandler(writer http.ResponseWriter, request *http.Request) {
	var item todoItem
	item.Title = request.FormValue("Title")
	database.addItem(item)
}

func deleteItemHandler(writer http.ResponseWriter, request *http.Request) {
	b, _ := ioutil.ReadAll(request.Body)
	var item todoItem
	json.Unmarshal(b, &item)
	database.deleteItem(item)
}

func findItemHandler(writer http.ResponseWriter, request *http.Request) {
	b, _ := ioutil.ReadAll(request.Body)
	var item todoItem
	json.Unmarshal(b, &item)

	err := database.findItem(item)

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
	} else {
		writer.WriteHeader(http.StatusAccepted)
	}
}

func main() {
	database = newDatabaseConnection()

	http.HandleFunc("/api/todos", getTodoListHandler)
	http.HandleFunc("/api/newItem", addItemHandler)
	http.HandleFunc("/api/deleteItem", deleteItemHandler)
	http.HandleFunc("/api/findItem", findItemHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {log.Fatal(err)}
}