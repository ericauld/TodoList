#!/bin/bash

main() {
    local prerequisite_packages="wget gcc yum"
    check_for_prerequisites $prerequisite_packages

    local install_list="go node yarn mysql foo bar"
    install_if_missing $install_list

    echo "Reached end of main"
}

check_for_prerequisites() {
    local prerequired_packages=$@
    local prerequired_missing_packages=$(which_are_missing $prerequired_packages)
    if [ ! -z "$prerequired_missing_packages" ]; then
        echo "Some prerequired packages were missing:"
        for pckg in $prerequired_missing_packages; do
            echo $pckg
        done
    fi
}

install_if_missing() {
    local install_list=$@
    local missing_packages=$(which_are_missing $install_list)
    if [ ! -z "$missing_packages" ]; then
        ask_for_permission_to_install $missing_packages
        install_packages $missing_packages
    fi
}

which_are_missing() {
    local necessary_packages=$@
    local missing_packages=""

    for pckg in $necessary_packages; do
        if ! shell_knows_the_command $pckg; then
            missing_packages+="$pckg "
        fi
    done
    echo $missing_packages
}

ask_for_permission_to_install() {
    local missing_packages=$@
    echo "The following packages were not found:"
    for pckg in $missing_packages; do
        echo $pckg
    done
    printf "Would you like the installer to install the missing packages? [y/N]"
    read yn
    case $yn in
    y)
        return
        ;;
    n)
        echo "Install script canceled because user declined to proceed."
        exit 0
        ;;
    *)
        echo "Invalid input. Install script will terminate."
        exit 1
        ;;
    esac
}

install_packages() {
    local packages_to_be_installed=$@
    for pckg in $packages_to_be_installed; do
        install_package $pckg
    done
}

shell_knows_the_command() {
    local command_name=$1
    command -v $command_name &>/dev/null
}

install_package() {
    local package_name=$1
    case $package_name in
    "mysql")
        install_mysql
        ;;
    "go")
        install_go
        ;;
    "node")
        configure_nodesource_repository
        ;;
    "yarn")
        install_yarn
        ;;
    "foo")
        install_foo
        ;;
    *)
        echo "Unknown package name \"$package_name\" given to install_package method"
        return 1
        ;;
    esac
}

install_mysql() {
    if shell_knows_the_command apt-get; then
        apt-get update
        apt-get install mysql-shell
    elif shell_knows_the_command wget; then
        sudo wget https://dev.mysql.com/get/mysql57-community-release-el7-11.noarch.rpm
        sudo yum localinstall mysql57-community-release-el7-11.noarch.rpm
        sudo yum install mysql-community-server
        sudo systemctl start mysqld.service
        sudo rm mysql57-community-release-el7-11.noarch.rpm
    else
        echo "Installation of mysql could not proceed, because neither \
\"apt-get\" nor \"wget\" commands were recognized."
    fi

    temp_password=$(sudo grep 'temporary password' /var/log/mysqld.log | awk '{print $NF}')

    echo -e "MySQL requires the creation of a password. This password must \n\
    --contain at least 1 numeric character\n\
    --contain at least 1 lowercase character\n\
    --contain at least 1 uppercase character\n\
    --contain at least 1 special (nonalphanumeric) character.\n\
    This password will be stored in a file \"password.txt\" which \
is already added to .gitignore."
    echo -n "Please enter a password for your MySQL database:"
    IFS= read -s password

    mysqladmin --user=root --password="$temp_password" password "$password"
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
    cat <<EOF >password.txt
$password
EOF

}

install_go() {
    if shell_knows_the_command wget; then
        sudo wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz
        sudo tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz
        sudo echo 'export PATH=/usr/local/go/bin:$PATH' >>~/.bash_profile
        source ~/.bash_profile
        if shell_knows_the_command go; then
            echo $(go version) successfully installed
        else
            echo There was a problem installing go
            exit
        fi
        rm go1.15.6.linux-amd64.tar.gz
    else
        echo "Could not install go because the \"wget\" package was not available"
    fi
}

configure_nodesource_repository() {
    curl --silent --location https://rpm.nodesource.com/setup_12.x | sudo bash -
}

install_yarn() {
    curl --silent --location https://dl.yarnpkg.com/rpm/yarn.repo | sudo tee /etc/yum.repos.d/yarn.repo
    sudo yum install yarn
}

install_foo() {
    echo "Installing foo..."
}

main
