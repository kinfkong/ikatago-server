#!/bin/bash
WEIGHTS_TO_INSTALL=$@
WEIGHT_SERVER_PORT=8022
WEIGHT_SERVER_HOST=106.13.32.114
WEIGHT_SERVER_USER=aistudio
WEIGHT_REMOTE_STORAGE_PATH=/home/aistudio/weights
WEIGHT_SERVER_KEY=./resources/weight_scripts/weight_server_key.pem
WEIGHT_SERVER_OPTIONS="-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i $WEIGHT_SERVER_KEY"
scp_weight() {
    FILE_PATH=$1
    WEIGHT_NAME=$2
    TARGET_MD5=$3
    SHOULD_UPDATE=0
    WEIGHT_REMOTE_FILE_PATH=$WEIGHT_REMOTE_STORAGE_PATH/$WEIGHT_NAME.bin.gz
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
        TARGET_DIR=$(dirname $FILE_PATH)
        mkdir -p $TARGET_DIR
        rm -rf $FILE_PATH.downloading
        scp $WEIGHT_SERVER_OPTIONS -P $WEIGHT_SERVER_PORT $WEIGHT_SERVER_USER@$WEIGHT_SERVER_HOST:$WEIGHT_REMOTE_FILE_PATH $FILE_PATH.downloading < /dev/null
        if [ $? -ne 0 ]
        then
            echo "failed to do scp for file: $WEIGHT_NAME, will try later"
            return -1
        fi
        mv $FILE_PATH.downloading $FILE_PATH
    fi
}

WEIGHT_STORAGE=./resources/extra_weights
mkdir -p ./resources/weight_scripts
mkdir -p $WEIGHT_STORAGE

wget https://ikatago-resources.oss-cn-beijing.aliyuncs.com/all/weight_server_key.pem -O $WEIGHT_SERVER_KEY
chmod 600 $WEIGHT_SERVER_KEY

for w in $WEIGHTS_TO_INSTALL
do
    echo "installing weight: $w, please wait..."
    scp_weight $WEIGHT_STORAGE/$w.bin.gz $w >/dev/null 2>&1 
    if [ $? -ne 0 ]
    then
        echo "failed to download weight $w"
        exit -1
    fi
done

mkdir -p ./work/data/weights
cp $WEIGHT_STORAGE/* ./work/data/weights/
rm -rf ./resources/weight_scripts
echo "Weights installed successfully!"