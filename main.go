package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	password, err := ioutil.ReadFile("./password.txt")
	if err != nil {log.Fatal(err)}

	_, err = sql.Open("mysql",
						fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/TodoList", password))
	if err != nil {log.Fatal(err)}
}