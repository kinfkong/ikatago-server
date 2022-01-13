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

PLATFORM_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhRW5jcnlwdEtleVByZWZpeCI6ImFsbC1zZWNyZXQiLCJpYXQiOjE2MDI4NjgxMzcsImV4cCI6MTYzMzAxNzU5OSwiYXVkIjoiYWxsIiwiaXNzIjoia2luZmtvbmcifQ.N7dSGORtxY5E0Qg-Dp0z_2EsbX0Icv6rHM-Daf0ZhPRBJ6ZcWU8Oiyna9flMU3EcVzQ8h0-AREIKsHvbbCrL3MkAIKNuGDG1JXt0gJB1fpyQ1WXYkfRn16nPNgwxxQAJ2wtrZImhLI4MWqYutFOxgzRY9iWTOLYXVA-DltFU89IsuVZ-hgTQ4oIHW6lrr8MwHnTYFpuTEOdUwYQT1pyFS4RHnWQLwupbT-zsa7EG_6wagFd_aAOLR08xQgY189YgJ5WVAKtFQjMKHdmfL0J5VmPqpLv8pTsbLKlKMXIMPXTyhyLxc5qHygVFyyZcEndUHNRu6vWy_z-wPXSac239VQ"


echo $USER_NAME":"$USER_PASSWORD > ./userlist.txt
chmod +x ./ikatago-server
./ikatago-server --platform all --token $PLATFORM_TOKEN

