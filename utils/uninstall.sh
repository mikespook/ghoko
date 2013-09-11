#!/bin/bash

usage() {
    echo "Usage: $0 -p"
    exit 1
}

PURGE=0

while getopts 'p' o &>> /dev/null; do
    case "$o" in
	p)
		PURGE=1;;
    *)
        usage;;
    esac
done

# check root
[ $EUID -ne 0 ] && echo '[ERROR] root needed' && exit 1

APP=ghoko

rm -f /usr/bin/$APP
update-rc.d -f $APP remove > /dev/null
rm -f /etc/init.d/$APP

if [ $PURGE -eq 1 ]; then
	rm -rf /usr/share/$APP
	rm -rf /etc/$APP.conf
fi

