#!/bin/bash

uninstall_mysql() {
    sudo yum remove mysql mysql-server
    sudo rm -fr /var/lib/mysql
    sudo rm -fr /etc/selinux/targeted/active/modules/100/mysql
    sudo rm -fr /usr/lib64/mysql
    sudo rm -fr /usr/share/mysql
    sudo rm -fr ~/mysql*
    sudo rm -f /var/log/mysql*
    #change to if doesnt exist look at datadir variable in  my.cnf change all these to if existx
    #make it a list
}

uninstall_go() {
    sudo rm -rvf /usr/local/go/
}

uninstall_mysql
uninstall_go
sudo rm -fr TodoList