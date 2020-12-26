package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type DatabaseConnection struct {
	db *sql.DB
}

func newDatabaseConnection() *DatabaseConnection {
	loginInfo := getLoginString2()
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
	rows, err := databaseConnection.db.Query(SQLQuery)
	if err != nil {log.Fatal(err)}
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

func getLoginString2() string {
	port := 3306
	databaseName := "TodoList"
	IPAddress := "127.0.0.1"
	username, password, _ := credentials()

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, IPAddress, port, databaseName)
}

func getLoginString() string {
	username := "root"
	port := 3306
	databaseName := "TodoList"
	IPAddress := "127.0.0.1"
	passwordAsByteSlice, err := ioutil.ReadFile("password.txt")
	password := string(passwordAsByteSlice)
	password = strings.TrimSuffix(password, "\n")
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, IPAddress, port, databaseName)
}

func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}