#!/bin/bash
mkdir ./tmp
cp -r ./config ./tmp
cd ./tmp
rm ./config/config.go ./config/frpc-natfrp.ini
zip -r config.zip ./config 
cp config.zip ../
cd -
rm -rf ./tmp