#!/bin/bash
rm -rf ./data && mkdir -p ./data
EN_PASS=abcde12345
openssl enc -in ./baidu-aistudio/katago -aes-256-cbc -pass pass:$EN_PASS > ./data/k
openssl enc -in ./baidu-aistudio/libstdc++.so.6.0.28 -aes-256-cbc -pass pass:$EN_PASS > ./data/lc
openssl enc -in ./baidu-aistudio/libzip.so.4.0.0 -aes-256-cbc -pass pass:$EN_PASS > ./data/lz
zip -r data.zip ./data