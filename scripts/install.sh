#!/bin/bash
PACKAGE=$1
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
        if [ $? -ne 0 ]
        then
            return -1
        fi
        mv $FILE_PATH.downloading $FILE_PATH
    fi
}
echo "Downloading package..."
update_file ./work.zip https://ikatago-resources.oss-cn-beijing.aliyuncs.com/all/ubuntu-18-cuda-10-2/work.zip 3dd48075a3a0d628aadb68641b3cd8da
if [ $? -ne 0 ]
then
    echo "Failed to download the package."
    exit -1
fi
echo "Download done"
echo "Installing..."
rm -rf work
unzip work.zip >/dev/null 2>&1
cd work
./install.sh >/dev/null 2>&1
echo "Install Successfully, now you can run the ikatago-server"