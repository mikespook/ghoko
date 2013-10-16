#!/bin/bash

[ $EUID -ne 0 ] && echo 'root needed' && exit 1

cp init.d_ghoko /etc/init.d/ghoko
cp ghoko.conf /etc/ghoko.conf
mkdir -p /usr/share/ghoko

update-rc.d -f ghoko remove
update-rc.d ghoko defaults
service ghoko restart
