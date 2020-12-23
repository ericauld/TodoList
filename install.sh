#!/bin/bash

install_yarn() {
    curl --silent --location https://dl.yarnpkg.com/rpm/yarn.repo | sudo tee /etc/yum.repos.d/yarn.repo
    curl --silent --location https://rpm.nodesource.com/setup_12.x | sudo bash -
    sudo yum install yarn
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

    echo "MySQL requires the creation of a password. This password must \n\
    --contain at least 1 numeric character\n\
    --contain at least 1 lowercase character\n\
    --contain at least 1 uppercase character\n\
    --contain at least 1 special (nonalphanumeric) character.\n\
    This password willl be stored in a file \"password.txt\" which\
    is already added to .gitignore."
    echo -n "Please enter a password for your MySQL database:"
    IFS= read -s password

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

    touch password.txt
    cat <<EOF > password.txt
$password
EOF

}

install_go() {
    if shellKnowsTheCommand wget; then
        sudo wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz
        sudo tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz
        sudo echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.bash_profile
        source ~/.bash_profile
        if shellKnowsTheCommand go; then
            echo $(go version) successfully installed
        else
            echo There was a problem installing go
            exit
        fi
    else
        echo "Could not install go because the \"wget\" package was not available"
    fi
}

shellKnowsTheCommand() {
    command_name=$1
    command -v $command_name &>/dev/null
}

main() {
    command_name=yarn
    if ! shellKnowsTheCommand $command_name; then
        while true; do
            read -p "$command_name was not found. Would you like to install it? [y/N]" yn
            case $yn in
            [Yy]*)
                install_yarn
                break
                ;;
            [Nn]*) exit ;;
            *) exit ;;
            esac
        done
    else
        echo I found $command_name on the system.
    fi  
        
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

    go get -u github.com/go-sql-driver/mysql
    yarn install
}
main
