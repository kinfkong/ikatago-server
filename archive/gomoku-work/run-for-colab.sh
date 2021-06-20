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

PLATFORM_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhRW5jcnlwdEtleVByZWZpeCI6ImpibGtqbGF4IiwiaWF0IjoxNjIyNzg4NTg5LCJleHAiOjE5NDg1NzkxOTksImF1ZCI6ImNvbGFiIiwiaXNzIjoia2luZmtvbmcifQ.WMPaTYJAuGx1QbUTrag5eX0e8pVU4eXCxoNlX4h2wrpOV3dMPSfi4boQvUkeAWreWsehNd9o7OxvdGpNQ0r8bIBLITVgoBDTGVTjxrJRrHCIgMa08HIohgwTjInW8SuPNZGFsKrUUnwAqCgS-6VDmc5TKd-t56DJyH6m3I0ELv26jjF7OzlhrSKlIz9HwYxh3OyU1qbsYaKQx74vs1ykacAvHJ4DQETxMmJPLpMOOmA9L7r26Qc8iFXcS5HEaDj-nZDUM471itIHT91QtgjPm9kdSVsO3k20MrOmerB0TM-gVxnEjEyjCfZGwdgGnbfYthBw96QbA6Mhwbf7ipXtlw"


echo $USER_NAME":"$USER_PASSWORD > ./userlist.txt
chmod +x ./change-frpc.sh 
./change-frpc.sh $USER_NAME
chmod +x ./ikatago-server
./ikatago-server --platform colab --token $PLATFORM_TOKEN

