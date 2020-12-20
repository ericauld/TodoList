package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type DatabaseConnection struct {
	db *sql.DB
}

func newDatabaseConnection() *DatabaseConnection {
	loginInfo := getLoginString()
	fmt.Println("Login string was " + loginInfo)
	fmt.Println("Got to just before opening db")
	db, err := sql.Open("mysql", loginInfo)
	fmt.Println("Got to just after opening db")
	err2 := db.Ping()
	if err != nil {log.Fatal(err)}
	fmt.Println("Err1 was nil")
	if err2 != nil {log.Fatal(err2)}
	fmt.Println("Err2 was nil")
	return &DatabaseConnection{db}
}

func (databaseConnection *DatabaseConnection) close() error {
	return databaseConnection.close()
}

func (databaseConnection *DatabaseConnection) ping() error {
	return databaseConnection.db.Ping()
}

func (databaseConnection *DatabaseConnection) findItem(item todoItem) error {
	nMatchingItems, err := databaseConnection.countItemsWhoseTitleIs(item.Title)
	if err != nil {return err}

	if nMatchingItems == 0 {
		return errors.New("no row with the given title was found")
	}
	return nil
}

func (databaseConnection *DatabaseConnection) addItem(item todoItem) error {
	SQLQuery, err := databaseConnection.db.Prepare(
		"INSERT INTO tasks(title) VALUES(?);")
	defer SQLQuery.Close()
	if err != nil {return err}
	_, err = SQLQuery.Exec(item.Title)
	return err
}

func (databaseConnection *DatabaseConnection) deleteItem(item todoItem) error {
	SQLQuery, err := databaseConnection.db.Prepare(
		"DELETE FROM tasks WHERE title=?;")
	defer SQLQuery.Close()
	if err != nil {return err}
	_, err = SQLQuery.Exec(item.Title)
	return err
}

func (databaseConnection *DatabaseConnection) countItemsWhoseTitleIs(itemTitle string) (int, error) {
	SQLQuery, err := databaseConnection.db.Prepare(
		"SELECT COUNT(*) FROM tasks WHERE title=?")
	defer SQLQuery.Close()
	if err != nil {
		return 0, err
	}
	var nMatchingRows int
	row := SQLQuery.QueryRow(itemTitle)
	row.Scan(&nMatchingRows)
	return nMatchingRows, nil
}

func (databaseConnection *DatabaseConnection) getTodoList() []todoItem {
	SQLQuery := "SELECT title FROM tasks;"
	fmt.Println("Got to just before calling query")
	rows, err := databaseConnection.db.Query(SQLQuery)
	fmt.Println("Got to just after calling query")
	if err != nil {log.Fatal(err)}
	fmt.Println("Got to line after log fatal err")
	defer rows.Close()

	var todoItems []todoItem
	for rows.Next() {
		var item todoItem
		err = rows.Scan(&item.Title)
		if err != nil {log.Fatal(err)}
		todoItems = append(todoItems, item)
	}
	return todoItems
}

func getLoginString() string {
	username := "root"
	port := 3306
	databaseName := "TodoList"
	IPAddress := "127.0.0.1"
	passwordAsByteSlice, err := ioutil.ReadFile("password.txt")
	password := string(passwordAsByteSlice)
	fmt.Println(password)
	password = strings.TrimSuffix(password, "\n")
	fmt.Println(password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v%v", password, password)

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, IPAddress, port, databaseName)
}

