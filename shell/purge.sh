#!/bin/bash

[ $EUID -ne 0 ] && echo 'root needed' && exit 1

service ghoko stop
update-rc.d -f ghoko remove

rm -rf /usr/share/ghoko

rm /etc/init.d/ghoko
rm /etc/ghoko.conf
