#!/bin/bash
UNAME=$1
CMD="s|\[kinfkong-ssh.*\]|[kinfkong-ssh-$UNAME]|g"
sed -i  -e "$CMD" ./config/frpc.txt
