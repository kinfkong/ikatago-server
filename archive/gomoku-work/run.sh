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

PLATFORM_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhRW5jcnlwdEtleVByZWZpeCI6ImxmZHNreCIsImlhdCI6MTYyMjc4ODU4OSwiZXhwIjoxOTQ4NTUwMzk5LCJhdWQiOiJhbGwiLCJpc3MiOiJraW5ma29uZyJ9.0doLI-nXqDlr0w0RRV2sBeEAZhS4mkrBg5YpjwqmnNubgnAulsCK9FrSLgDJGMFSh_lW2XM1SaZWQ6A9-GnsS0t6rY_41_xyJ7eE35MFI3crOkT8eJIMzGapuxpPsNBjPHYB3CgduMZopmoQqLgHNMpXwQE-Ed43wCxU-_zMYJHTyHCf_Anbs0NEo9P8H8ocvLQ5V5GppXmQQa3-PgO9FQP1HiL4iEa0W2F9AH--YE1V6Hd2IoeX_1i4RSzdDB0dGDcjsA59BTxCcHGbnwx30IYaaDfIuGVw7AmFfUR3qQJln3XUTqV_EPezS-xHD1v1-BvNocxvCI9evPe5ftAEDA"


echo $USER_NAME":"$USER_PASSWORD > ./userlist.txt

chmod +x ./ikatago-server
./ikatago-server --platform all --token $PLATFORM_TOKEN

