#!/bin/bash

main() {
    prerequisites="git yarn gcc node wget yum"
    make_sure_system_already_has $prerequisites

    local packages="go mysql"
    install_packages $(missing $packages)

    install_drivers_and_dependencies

    setup_mysql_database
}

make_sure_system_already_has() {
    local prerequisites=$@
    local missing_prerequisites=$(missing $prerequisites)
    if [ ! -z "$missing_prerequisites" ]; then
        tell_user_to_install $missing_prerequisites
        exit 0
    fi
}

install_packages() {
    packages=$@
    for pckg in $packages; do
        eval install_$pckg
    done 
}

install_drivers_and_dependencies() {
    install_go_mysql_driver
    install_yarn_dependencies
}

setup_mysql_database() {
    local temp_password=$(get_temp_mysql_password)
    ask_user_to_type_in_a_password
    read -rs user_password
    update_password $temp_password $user_password
    save_password_to_gitignored_file $user_password
    create_database $user_password
}

missing() {
    local packages=$@
    local missing_packages=""

    for pckg in $packages; do
        if system_cannot_find $pckg; then
            missing_packages+="$pckg "
        fi
    done
    echo $missing_packages
}

tell_user_to_install() {
    local missing_packages=$@
    echo "In order to install TodoList, the following must already be installed:"
    for prerequisite in $prerequisites; do
        echo -e "   * "$prerequisite
    done
    echo "But the following were not found on your system:"
    for pckg in $missing_packages; do
        echo -e "   * "$pckg
    done
}

install_mysql() {
    sudo wget https://dev.mysql.com/get/mysql57-community-release-el7-11.noarch.rpm
    sudo yum localinstall mysql57-community-release-el7-11.noarch.rpm
    sudo yum install mysql-community-server
    sudo systemctl start mysqld.service
    sudo rm mysql57-community-release-el7-11.noarch.rpm
}

install_go() {
    sudo wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz
    sudo echo 'export PATH=/usr/local/go/bin:$PATH' >>~/.bash_profile
    source ~/.bash_profile
    rm go1.15.6.linux-amd64.tar.gz
}

install_go_mysql_driver() {
    go get -u github.com/go-sql-driver/mysql
}

install_yarn_dependencies() {
    yarn install
}

get_temp_mysql_password() {
    local temp_password=$(sudo grep 'temporary password' /var/log/mysqld.log | awk '{print $NF}')
    echo $temp_password
}

ask_user_to_type_in_a_password() {
    echo -e "MySQL requires the creation of a password. This password must \n\
    --contain at least 1 numeric character\n\
    --contain at least 1 lowercase character\n\
    --contain at least 1 uppercase character\n\
    --contain at least 1 special (nonalphanumeric) character.\n\
    This password will be stored in a file \"password.txt\" which \
is already added to .gitignore."
    echo -n "Please enter a password for your MySQL database:"
}

update_password() {
    local temp_password=$1
    local password=$2
    mysqladmin --user=root --password="$temp_password" password "$password"
}

create_database() {
    local password=$1
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

save_password_to_gitignored_file() {
    local password=$1
    touch password.txt
    cat <<EOF >password.txt
$password
EOF

}

system_cannot_find() {
    local pckg=$1
    if command -v $pckg &>/dev/null; then
        return 1
    else
        return 0
    fi
}

main