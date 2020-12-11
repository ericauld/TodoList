# Todo List

A simple todo list implemented with React, MySQL, and a Go backend. This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

## Requirements

### React

First, you need to [install](https://classic.yarnpkg.com/en/docs/install/#mac-stable) Yarn package manager. Then you can navigate to the project directory in the terminal, and run the command `yarn install`. It should create a folder `<project-directory>/node_modules`, where the necessary packages are installed. There is a `yarn.lock` file in the repository, which instructs Yarn which packages are necessary.

### Go

To run the server, you should [have Go installed](http://golang.org/doc/install.html) on your machine. You should also have the [Go MySQL driver](https://github.com/go-sql-driver/mysql) installed.  

### MySQL

Currently the app is set up to run off a MySQL database hosted on your local machine. Therefore, you should create a database to hold your todo items. Currently the only requirement of the database is that is has a table called `tasks` which has a field called `title`, which is a `VARCHAR(255)`.

Running the command

    go run <project-dir>/main.go <project-dir>/database.go

will log you into your database and start the server. You may care to look at the function `database.go::getLoginString`, where the login string is formatted. There you will see that the username is currently set to `root`, the port is set to `3306`, and the database title is set to `TodoList`. All of these are adjustable, of course. The password for the database is read in from a file called `password.txt`, which you should create in the project directory. It should have nothing but the password in it (and no newline afterward). It has already been placed in `.gitignore` for you.

## Operation

In the project directory, you can run 

    go run main.go database.go

to start the server. Then, in a separate terminal process, simply type 

    yarn start

This will start up the frontend and automatically open the site in a new browser window. Currently the backend operates on port `8080` and the frontend on port `3000`. So the site will be hosted on `localhost:3000`.
