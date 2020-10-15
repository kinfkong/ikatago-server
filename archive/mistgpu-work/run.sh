#!/bin/bash

USER_NAME=$1
USER_PASSWORD=$2
if [ "$USER_NAME" == "" ]
then
	echo "Usage: $0 USER_NAME USER_PASSWORD"
	exit -1
fi
if [ "$USER_PASSWORD" == "" ]
then
	echo "Usage: $0 USER_NAME USER_PASSWORD"
	exit -1
fi

PLATFORM_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhRW5jcnlwdEtleVByZWZpeCI6ImFpc3R1ZGlvLXNlY3JldCIsImlhdCI6MTU5Njk1MTgwMywiZXhwIjoxNjMzMDQ2Mzk5LCJhdWQiOiJhaXN0dWRpbyIsImlzcyI6ImtpbmZrb25nIn0.US9KvL96UkPsbxyhrABa4IY2FXeK4-JNT_zspey5kG2yRhp3fKBZw7Um0W3wBXzGrGc6-DvPkDRvkgF4B_aOAMT9aYDBaI3PpThIL10tmGE7nCtWb-7yDCYd8whTlaiQMnhT7uQCKyCQ4zTxB9DaVbDz0chxmmpEaHUBBJely5YTNb7OMgDNwIVmx0VYckOY7miPgQQBJmYzFJeut8eUBiDsCvBdjs0rvhiKddyIvqXMs73aNFPvNcuLqErcq3UytoLrr9tB5Wa-eo1DeIAr3RwJJuID5AoT-6zMBETIPJnKGeDNCYZcaUumDuR5bGc56_vbFVfJQxgYmSe8RIFJ8A"


# check upgrade and do auto upgrade
wget -q https://ikatago-resources.oss-cn-beijing.aliyuncs.com/upgrade.sh -O ./upgrade.sh
chmod +x ./upgrade.sh
./upgrade.sh

echo $USER_NAME":"$USER_PASSWORD > ./userlist.txt
chmod +x ./change-frpc.sh 
./change-frpc.sh $USER_NAME
chmod +x ./ikatago-server
./ikatago-server --platform aistudio --token $PLATFORM_TOKEN

