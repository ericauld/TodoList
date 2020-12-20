#!/bin/bash

testIfInstalledAndPromptForInstallation() {
    command_name=$1
    if ! command -v $command_name &>/dev/null; then
        while true; do
            read -p "$command_name was not found. Would you like to install it? [y/N]" yn
            case $yn in
            [Yy]*)
                echo "Here I go! I'm installing some shit"
                break
                ;;
            [Nn]*) return ;;
            *) return ;;
            esac
        done
    else
        echo "$command_name confirmed installed"
    fi
}

# vercomp is used from
# https://stackoverflow.com/questions/4023830/how-to-compare-two-strings-in-dot-separated-version-format-in-bash
vercomp() {
    if [[ $1 == $2 ]]; then
        return 0
    fi
    local IFS=.
    local i ver1=($1) ver2=($2)
    # fill empty fields in ver1 with zeros
    for ((i = ${#ver1[@]}; i < ${#ver2[@]}; i++)); do
        ver1[i]=0
    done
    for ((i = 0; i < ${#ver1[@]}; i++)); do
        if [[ -z ${ver2[i]} ]]; then
            # fill empty fields in ver2 with zeros
            ver2[i]=0
        fi
        if ((10#${ver1[i]} > 10#${ver2[i]})); then
            return 1
        fi
        if ((10#${ver1[i]} < 10#${ver2[i]})); then
            return 2
        fi
    done
    return 0
}

install_mysql() {
    if shellKnowsTheCommand apt-get; then
        apt-get update
        apt-get install mysql-shell
    elif shellKnowsTheCommand wget; then
        sudo wget https://dev.mysql.com/get/mysql57-community-release-el7-11.noarch.rpm
        sudo yum localinstall mysql57-community-release-el7-11.noarch.rpm
        sudo yum install mysql-community-server
        sudo systemctl start mysqld.service

    else
        echo Hi
    fi

    password="MyNewPass4%"
    temp_password=$(sudo grep 'temporary password' /var/log/mysqld.log | awk '{print $NF}')  
    mysqladmin --user=root --password="$temp_password" password "$password";
    mysql -uroot -p"$password" -e "CREATE DATABASE TodoList;"
    mysql -uroot -p"$password" -e "USE TodoList; CREATE TABLE IF NOT EXISTS tasks (
    task_id INT AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    start_date DATE,
    due_date DATE,
    priority TINYINT NOT NULL DEFAULT 3,
    description TEXT,
    PRIMARY KEY (task_id));"
    mysql -uroot -p"$password" -e "USE TodoList; INSERT INTO tasks(title,priority)
VALUES('Learn MySQL INSERT Statement',1);"



}

install_go() {
    echo install_go was called
    if shellKnowsTheCommand wget; then
        sudo wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz
        sudo tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz
        sudo echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.bash_profile
        source ~/.bash_profile
        if shellKnowsTheCommand go; then
            echo go installed with version $(go version)
        else
            echo There was a problem installing go
            exit
        fi
    else
        echo "Could not install go because wget wasn't recognized"
    fi
}

install_git() {
    echo "Install git was called."
    sudo yum update -y
    sudo yum install git -y
}

shellKnowsTheCommand() {
    command_name=$1
    command -v $command_name &>/dev/null
}

main() {
    # testIfInstalledAndPromptForInstallation yarn
    # testIfInstalledAndPromptForInstallation mysql
    # testIfInstalledAndPromptForInstallation go

    sudo yum groupinstall "Development Tools"

    command_name=mysql
    if ! shellKnowsTheCommand $command_name; then
        while true; do
            read -p "$command_name was not found. Would you like to install it? [y/N]" yn
            case $yn in
            [Yy]*)
                install_mysql
                break
                ;;
            [Nn]*) exit ;;
            *) exit ;;
            esac
        done
    else
        echo I found $command_name on the system.
    fi

    command_name=go
    if ! shellKnowsTheCommand $command_name; then
        while true; do
            read -p "$command_name was not found. Would you like to install it? [y/N]" yn
            case $yn in
            [Yy]*)
                install_go
                break
                ;;
            [Nn]*) exit ;;
            *) exit ;;
            esac
        done
    else
        echo I found $command_name on the system.
    fi

    command_name=git
    if ! shellKnowsTheCommand $command_name; then
        while true; do
            read -p "$command_name was not found. Would you like to install it? [y/N]" yn
            case $yn in
            [Yy]*)
                install_git
                break
                ;;
            [Nn]*) exit ;;
            *) exit ;;
            esac
        done
    else
        echo I found $command_name on the system.
    fi

    mkdir -p ~/go/src/github.com/ericauld ; cd ~/go/src/github.com/ericauld
    git clone https://github.com/ericauld/TodoList.git
    cd ~/go/src/github.com/ericauld/TodoList/

    go get -u github.com/go-sql-driver/mysql
    yarn install
}
main
