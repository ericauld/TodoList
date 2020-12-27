#!/bin/bash

main() {
    users_os=$(get_users_os)
    make_sure_script_supports_users_os

    prerequisites="git yarn gcc node "
    prerequisites+=$(os_specific_prerequisites_for $users_os)
    make_sure_system_already_has_prerequisites

    local packages_to_install="go mysql"
    local missing_packages=$(missing $packages_to_install)
    get_users_permission_to_install $missing_packages
    install_packages $missing_packages

    install_drivers_and_dependencies

    setup_mysql_database
}

get_users_os() {
    case $OSTYPE in
    ( darwin* )
        echo mac;;
    ( linux* )
        echo $(get_linux_distribution);;
    ( * )
        echo "unknown_os";;
    esac
}

make_sure_script_supports_users_os() {
    local supported_os_types="amazon_linux mac"
    if isnt_a_member_of $users_os "$supported_os_types"; then
        echo "The script couldn't discern your OS as being one of the \
supported types. Currently supported are "
        for os_type in $supported_os_types; do
            echo $os_type
        done
    fi
}

os_specific_prerequisites_for() {
    local users_os=$1
    if [ $users_os = "amazon_linux" ]; then
        echo "yum wget"
    elif [ $users_os = "mac" ]; then
        echo "curl mac stuff"
    fi
}

make_sure_system_already_has_prerequisites() {
    local missing_prerequisites=$(missing $prerequisites)
    if [ ! -z "$missing_prerequisites" ]; then
        tell_user_to_install $missing_prerequisites
        exit 0
    fi
}

install_packages() {
    local packages=$@
    for pckg in $packages; do
        eval install_"$pckg"_"$users_os"
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
    save_password_to_gitignored_file $user_password
    setup_login_path $temp_password $user_password
    create_database
}

get_linux_distribution() {
    local linux_distribution="$(grep '^NAME' /etc/os-release | sed 's/.*=//' | sed "s/^\(\"\)\(.*\)\1\$/\2/g")"
    if [ "$linux_distribution" = "Amazon Linux" ]; then
        echo amazon_linux
    else 
        echo unknown_linux_distribution
    fi
}

isnt_a_member_of() {
    local element=$1
    local list=$2
    [[ $list =~ (^|[[:space:]])$element($|[[:space:]]) ]] && return 1 || return 0
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

get_users_permission_to_install() {
    local missing_packages=$@
    echo "May the TodoList installer install the following packages on your machine?"
    for pckg in $missing_packages; do
        echo $pckg
    done
    read -p "[y/N] " -n 1 -r
    echo    # move to a new line
    if [[ ! $REPLY =~ ^[Yy]$ ]]
    then
        echo "User declined to install packages. Install script exiting"
        exit 1
    fi
}

tell_user_to_install() {
    local missing_packages=$@
    echo "In order to install TodoList on $users_os, the following must already be installed:"
    for prerequisite in $prerequisites; do
        echo -e "   * "$prerequisite
    done
    echo "But the following were not found on your system:"
    for pckg in $missing_packages; do
        echo -e "   * "$pckg
    done
}

install_mysql_amazon_linux() {
    sudo wget https://dev.mysql.com/get/mysql57-community-release-el7-11.noarch.rpm
    sudo yum localinstall mysql57-community-release-el7-11.noarch.rpm
    sudo yum install mysql-community-server
    sudo systemctl start mysqld.service
    sudo rm -f mysql57-community-release-el7-11.noarch.rpm
}

install_go_amazon_linux() {
    sudo wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz
    sudo echo 'export PATH=/usr/local/go/bin:$PATH' >>~/.bash_profile
    source ~/.bash_profile
    sudo rm -f go1.15.6.linux-amd64.tar.gz
}

install_go_mac() {
    curl -O "https://golang.org/dl/go1.15.6.darwin-amd64.pkg"
    ./go1.15.6.darwin-amd64.pkg
}

install_mysql_mac() {
    local name_of_download_file="mysql_installer"
    local path_to_todolist="~/go/src/github.com/ericauld/TodoList"
    curl -o $name_of_download_file "https://dev.mysql.com/downloads/file/?id=499570"
    sudo groupadd mysql
    sudo useradd -r -g mysql -s /bin/false mysql
    cd /usr/local
    sudo tar zxvf "$path_to_todolist"/"$name_of_download_file".tar.gz
    sudo ln -s "$path_to_todolist"/"$name_of_download_file" mysql
    cd mysql
    sudo mkdir mysql-files
    sudo chown mysql:mysql mysql-files
    sudo chmod 750 mysql-files
    sudo bin/mysqld --initialize --user=mysql
    sudo bin/mysql_ssl_rsa_setup
    #What's going on here?
    sudo bin/mysqld_safe --user=mysql &
    # Next command is optional
    sudo cp support-files/mysql.server /etc/init.d/mysql.server
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
After you type in your password you will be asked to type it once more\
in order to set up your login path."
    echo -n "Please enter a password for your MySQL database:"
}

setup_login_path() {
    local temp_password=$1
    local user_password=$2
    mysqladmin --user=root --password="$temp_password" password "$user_password"
    mysql_config_editor set --login-path=local --host=localhost --user=root --password
}

create_database() {
    mysql --login-path=local -e "CREATE DATABASE TodoList;"
    mysql --login-path=local -e "USE TodoList; CREATE TABLE IF NOT EXISTS tasks (
    task_id INT AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    start_date DATE,
    due_date DATE,
    priority TINYINT NOT NULL DEFAULT 3,
    description TEXT,
    PRIMARY KEY (task_id));"
    mysql --login-path=local -e "USE TodoList; INSERT INTO tasks(title,priority)
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