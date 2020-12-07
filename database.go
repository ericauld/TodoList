package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"fmt"
)

type DatabaseConnection struct {
	db *sql.DB
}

func openDatabaseConnection() *DatabaseConnection {
	loginInfo := getLoginString()
	db, err := sql.Open("mysql", loginInfo)
	if err != nil {log.Fatal(err)}
	return &DatabaseConnection{db}
}

func (databaseConnection *DatabaseConnection) ping() error {
	return databaseConnection.db.Ping()
}

func (databaseConnection *DatabaseConnection) getTodoList() []string {
	SQLQuery := "SELECT title FROM tasks;"
	rows, err := databaseConnection.db.Query(SQLQuery)
	if err != nil {log.Fatal(err)}
	defer rows.Close()

	var todoItemTitles []string
	for rows.Next() {
		var todoItemTitle string
		err = rows.Scan(&todoItemTitle)
		if err != nil {log.Fatal(err)}
		todoItemTitles = append(todoItemTitles, todoItemTitle)
	}
	return todoItemTitles
}

func getLoginString() string {
	username := "root"
	port := 3306
	databaseName := "TodoList"
	IPAddress := "127.0.0.1"
	password, err := ioutil.ReadFile("password.txt")
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, IPAddress, port, databaseName)
}

