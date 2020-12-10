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

	SQLQuery := "SELECT title FROM tasks;"
	rows, err := database.db.Query(SQLQuery)
	if err != nil {log.Fatal(err)}
	defer rows.Close()

	var todoItems []todoItem
	for rows.Next() {
		var item todoItem
		err = rows.Scan(&item.Title)
		if err != nil {log.Fatal(err)}
		todoItems = append(todoItems, item)
	}

	json.NewEncoder(writer).Encode(todoItems)
}

func addItemHandler(writer http.ResponseWriter, request *http.Request) {
	var item todoItem
	item.Title = request.FormValue("Title")

	SQLQuery, err := database.db.Prepare("INSERT INTO tasks(title) VALUES(?);")
	defer SQLQuery.Close()
	if err != nil {log.Fatal(err)}
	SQLQuery.Exec(item.Title)
}

func deleteItemHandler(writer http.ResponseWriter, request *http.Request) {
	b, _ := ioutil.ReadAll(request.Body)
	var item todoItem
	json.Unmarshal(b, &item)

	SQLQuery, err := database.db.Prepare("DELETE FROM tasks WHERE title=?")
	defer SQLQuery.Close()
	if err != nil {log.Fatal(err)}
	SQLQuery.Exec(item.Title)
}

func findItemHandler(writer http.ResponseWriter, request *http.Request) {
	b, _ := ioutil.ReadAll(request.Body)
	var item todoItem
	json.Unmarshal(b, &item)

	SQLQuery, err := database.db.Prepare("SELECT COUNT(*) FROM tasks WHERE title=?")
	defer SQLQuery.Close()
	if err != nil {log.Fatal(err)}
	var nMatchingRows int
	row := SQLQuery.QueryRow(item.Title)
	row.Scan(&nMatchingRows)

	if nMatchingRows == 0 {
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