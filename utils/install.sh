#!/bin/bash

usage() {
    echo "Usage: $0"
    exit 1
}

while getopts '' o &>> /dev/null; do
    case "$o" in
    *)
        usage;;
    esac
done

# check root
[ $EUID -ne 0 ] && echo '[ERROR] root needed' && exit 1

APP=ghoko

# load config
base=`dirname -- "$0"`
. $base/etc/$APP.conf

# install dependencies
apt-get -y -q -q install liblua5.1-0
echo "[INFO] dependencies installed"

# deploy config
for f in $base/etc/* ; do
	t=/etc/`basename $f`
	if [ ! -L $t ]; then
		cp -Rf $f /etc
	fi
done
echo "[INFO] configrations deployed"

# deploy script
mkdir -p $SCRIPT
cp -f $base/usr/share/$APP/* $SCRIPT
echo "[INFO] scripts deployed"

# deploy service
cp -f $base/usr/bin/$APP /usr/bin/
for s in $base/etc/init.d/*; do
    s=`basename $s`
    update-rc.d -f $s remove > /dev/null
    update-rc.d $s defaults > /dev/null
    service $s restart > /dev/null
done
echo "[INFO] service deployed"
