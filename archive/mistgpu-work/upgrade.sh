#!/bin/bash
update_file() {
    FILE_PATH=$1
    FILE_URL=$2
    TARGET_MD5=$3
    SHOULD_UPDATE=0
    if [ ! -f $FILE_PATH ]
    then
        SHOULD_UPDATE=1
    fi
    if [ $SHOULD_UPDATE -eq 0 ] && [ "$TARGET_MD5" != "" ]
    then
        #check md5
        MD5=`md5sum $FILE_PATH | cut -d' ' -f1`
        if [ "$MD5" != "$TARGET_MD5" ]
        then
            SHOULD_UPDATE=1
        fi
    fi
    if [ $SHOULD_UPDATE -eq 1 ]
    then
        rm -rf $FILE_PATH.downloading
        wget -q $FILE_URL -O $FILE_PATH.downloading -o /dev/null
        mv $FILE_PATH.downloading $FILE_PATH
    fi
}

# update_file ./data/weights/40b-large.bin.gz https://ikatago-resources.oss-cn-beijing.aliyuncs.com/40b384.bin.gz & 
update_file ./ikatago-server https://ikatago-resources.oss-cn-beijing.aliyuncs.com/ikatago-server 6489dcd1de2ea4975aa0d3db212b2e7e



