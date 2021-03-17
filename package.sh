#!/bin/bash
SVR_VERSION=1.5.1
BIN_DIR=./bin
BINARY_NAME=ikatago-server
rm -rf $BIN_DIR && mkdir $BIN_DIR


# Mac OSX
OUTPUT_PATH=$BINARY_NAME-$SVR_VERSION-mac-osx
mkdir -p $BIN_DIR/$OUTPUT_PATH
GOOS=darwin GOARCH=amd64 go build -o $BIN_DIR/$OUTPUT_PATH/$BINARY_NAME
cd $BIN_DIR
zip -r $OUTPUT_PATH.zip $OUTPUT_PATH
cd - >/dev/null

# linux
OUTPUT_PATH=$BINARY_NAME-$SVR_VERSION-linux
mkdir -p $BIN_DIR/$OUTPUT_PATH
GOOS=linux GOARCH=amd64 go build -o $BIN_DIR/$OUTPUT_PATH/$BINARY_NAME
cd $BIN_DIR
zip -r $OUTPUT_PATH.zip $OUTPUT_PATH
cd - >/dev/null

# windows 64bit
OUTPUT_PATH=$BINARY_NAME-$SVR_VERSION-win64
mkdir -p $BIN_DIR/$OUTPUT_PATH
GOOS=windows GOARCH=amd64 go build -o $BIN_DIR/$OUTPUT_PATH/$BINARY_NAME.exe
cd $BIN_DIR
zip -r $OUTPUT_PATH.zip $OUTPUT_PATH
cd - >/dev/null

# windows 32bit
OUTPUT_PATH=$BINARY_NAME-$SVR_VERSION-win32
mkdir -p $BIN_DIR/$OUTPUT_PATH
GOOS=windows GOARCH=386 go build -o $BIN_DIR/$OUTPUT_PATH/$BINARY_NAME.exe
cd $BIN_DIR
zip -r $OUTPUT_PATH.zip $OUTPUT_PATH
cd - >/dev/null


# raspbian pi
OUTPUT_PATH=$BINARY_NAME-$SVR_VERSION-raspbian
mkdir -p $BIN_DIR/$OUTPUT_PATH
GOOS=linux GOARCH=arm go build -o $BIN_DIR/$OUTPUT_PATH/$BINARY_NAME
cd $BIN_DIR
zip -r $OUTPUT_PATH.zip $OUTPUT_PATH
cd - >/dev/null

