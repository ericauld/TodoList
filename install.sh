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
        exit
    else
        echo "$command_name confirmed installed"
    fi
}

# vercomp is used from 
# https://stackoverflow.com/questions/4023830/how-to-compare-two-strings-in-dot-separated-version-format-in-bash
vercomp () {
    if [[ $1 == $2 ]]
    then
        return 0
    fi
    local IFS=.
    local i ver1=($1) ver2=($2)
    # fill empty fields in ver1 with zeros
    for ((i=${#ver1[@]}; i<${#ver2[@]}; i++))
    do
        ver1[i]=0
    done
    for ((i=0; i<${#ver1[@]}; i++))
    do
        if [[ -z ${ver2[i]} ]]
        then
            # fill empty fields in ver2 with zeros
            ver2[i]=0
        fi
        if ((10#${ver1[i]} > 10#${ver2[i]}))
        then
            return 1
        fi
        if ((10#${ver1[i]} < 10#${ver2[i]}))
        then
            return 2
        fi
    done
    return 0
}

main() {
    testIfInstalledAndPromptForInstallation mysql
    testIfInstalledAndPromptForInstallation yarn
    testIfInstalledAndPromptForInstallation go
}

main
