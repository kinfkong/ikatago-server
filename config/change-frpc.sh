#!/bin/bash
UNAME=$1
CMD="s|\[ssh.*\]|[ssh-$UNAME]|g"
sed -i  -e "$CMD" ./config/frpc.ini