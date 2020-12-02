package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
)

type todoItem struct {
	Title string
}

var db *sql.DB

func getTodoListEndpoint(writer http.ResponseWriter, request *http.Request) {
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

func addItemEndpoint(writer http.ResponseWriter, request *http.Request) {
	var item todoItem
	item.Title = request.FormValue("Title")

	SQLQuery, err := db.Prepare("INSERT INTO tasks(title) VALUES(?);")
	if err != nil {log.Fatal(err)}
	SQLQuery.Exec(item.Title)
}

func main() {
	password, err := ioutil.ReadFile("./password.txt")
	if err != nil {log.Fatal(err)}

	db, err = sql.Open("mysql",
						fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/TodoList", password))
	if err != nil {log.Fatal(err)}
	defer db.Close()

	http.HandleFunc("/api/todos", getTodoListEndpoint)
	http.HandleFunc("/api/newItem", addItemEndpoint)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {log.Fatal(err)}
}