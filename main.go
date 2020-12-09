package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
)

type todoItem struct {
	Title string
}

var db *sql.DB

var (
	database *DatabaseConnection
)

func getTodoListHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	SQLQuery := "SELECT title FROM tasks;"
	rows, err := db.Query(SQLQuery)
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

	SQLQuery, err := db.Prepare("INSERT INTO tasks(title) VALUES(?);")
	if err != nil {log.Fatal(err)}
	SQLQuery.Exec(item.Title)
}

func deleteItemHandler(writer http.ResponseWriter, request *http.Request) {
	b, _ := ioutil.ReadAll(request.Body)
	var item todoItem
	json.Unmarshal(b, &item)

	SQLQuery, err := db.Prepare("DELETE FROM tasks WHERE title=?")
	if err != nil {log.Fatal(err)}
	SQLQuery.Exec(item.Title)
}

func findItemHandler(writer http.ResponseWriter, request *http.Request) {
	b, _ := ioutil.ReadAll(request.Body)
	var item todoItem
	json.Unmarshal(b, &item)

	SQLQuery, err := db.Prepare("SELECT COUNT(*) FROM tasks WHERE title=?")
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
	password, err := ioutil.ReadFile("./password.txt")
	if err != nil {log.Fatal(err)}

	db, err = sql.Open("mysql",
						fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/TodoList", password))
	if err != nil {log.Fatal(err)}
	defer db.Close()

	http.HandleFunc("/api/todos", getTodoListHandler)
	http.HandleFunc("/api/newItem", addItemHandler)
	http.HandleFunc("/api/deleteItem", deleteItemHandler)
	http.HandleFunc("/api/findItem", findItemHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {log.Fatal(err)}
}