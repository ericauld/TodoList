package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

type DatabaseConnection struct {
	db *sql.DB
}

func newDatabaseConnection() *DatabaseConnection {
	loginInfo := getLoginString()
	db, err := sql.Open("mysql", loginInfo)
	if err != nil {log.Fatal(err)}
	return &DatabaseConnection{db}
}

func (databaseConnection *DatabaseConnection) close() error {
	return databaseConnection.close()
}

func (databaseConnection *DatabaseConnection) ping() error {
	return databaseConnection.db.Ping()
}

func (databaseConnection *DatabaseConnection) findItem(itemTitle string) error {
	nMatchingItems, err := databaseConnection.countItemsWhoseTitleIs(itemTitle)
	if err != nil {return err}

	if nMatchingItems == 0 {
		return errors.New("no row with the given title was found")
	}
	return nil
}

func (databaseConnection *DatabaseConnection) addItem(itemTitle string) error {
	SQLQuery, err := databaseConnection.db.Prepare(
		"INSERT INTO tasks(title) VALUES(?);")
	defer SQLQuery.Close()
	if err != nil {return err}
	_, err = SQLQuery.Exec(itemTitle)
	return err
}

func (databaseConnection *DatabaseConnection) deleteItem(itemTitle string) error {
	SQLQuery, err := databaseConnection.db.Prepare(
		"DELETE FROM tasks WHERE title=?;")
	defer SQLQuery.Close()
	if err != nil {return err}
	_, err = SQLQuery.Exec(itemTitle)
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

