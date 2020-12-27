
A simple todo list implemented with React, MySQL, and a Go backend. This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

## Requirements
So far the app has only been tested on Amazon Linux. It is in the process of being adapted for MacOS.

Before you install the app, you must have [Yarn](https://classic.yarnpkg.com/en/docs/install) and [Node](https://nodejs.org/en/download/package-manager). You may care to install Node with a package manager, as described [here](https://classic.yarnpkg.com/en/docs/install/#centos-stable) for Amazon Linux/Red Hat Linux. Also required for the app are gcc and git. 

If you are on Amazon Linux, you must also have `yum` and `wget` installed (but you probably already do).

If you use MacOS you will need `curl`. 

## Installation

Once you have installed the prerequisites, create a folder to install the app in. In order for Go packages to integrate smoothly, I recommend you use the directory `~/go/src/github.com/ericauld`, for instance with the terminal command

`mkdir -p ~/go/src/github.com/ericauld`

Then clone the repository into your new folder, for instance with the terminal commands

`cd ~/go/src/github.com/ericauld`

`git clone https://github.com/ericauld/TodoList.git`

Then navigate to the repository and run the `install.sh` script, as with the commands

`cd TodoList`

`./install.sh`

### What Installation Does

The installer will check that you have the required programs before installing (by simply checking that the shell knows the commands `git`, `gcc`, and so forth). If you do, it will check whether you have Go and MySQL installed (by the same method). If either of them are missing, it will ask you if you want to install them.

Then the installer will install `go-mysql-driver`. And it will run the command `yarn install`, which looks at the `yarn.lock` file and installs the dependencies it sees there.

Finally, the installer will configure MySQL and set up an appropriate database for the app. It will find the temporary password that was created for you when MySQL was installed, and prompt you to replace it with a password of your own. Then it will setup a login path for you on MySQL, so you don't have to keep putting in your password. Then it will create a database for you, and insert a dummy item into it.

## Operation

In the project directory, you can run 

    go run main.go database.go

to start the server. Then, in a separate terminal process, simply type 

    yarn start

This will start up the frontend and automatically open the site in a new browser window. Currently the backend operates on port `8080` and the frontend on port `3000`. So the site will be hosted on `localhost:3000`.

You may also wish to run one of the test scripts that are provided, for instance `testbackend`. Be careful if you do, since these will terminate any process on port `8080` that they find. 
