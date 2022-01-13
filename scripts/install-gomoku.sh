#!/bin/bash
DATA_DOWNLOAD_URL="http://8.130.167.86:8080/api/ikatago/resources/download/gomoku-all"

OS_NAME=$(cat /etc/os-release | grep "PRETTY_NAME" | sed 's/PRETTY_NAME="\(.*\)"/\1/g' | tr '[:upper:]' '[:lower:]')

GPU_NAME=$(nvidia-smi -q | grep "Product Name" | head -n 1 | cut -d":" -f2 | xargs)
if [[ "$GPU_NAME" == *"2080 Ti"* ]]
then
    GPU_NAME="2080Ti"
elif [[ "$GPU_NAME" == *"V100"* ]]
then
    GPU_NAME="v100"
elif [[ "$GPU_NAME" == *"RTX 3090"* ]]
then
    GPU_NAME="3090"
elif [[ "$GPU_NAME" == *"A100"* ]]
then
    GPU_NAME="A100"
fi

if [ -f /usr/local/cuda/version.txt  ]
then
    CUDA_VERSION=$(cat /usr/local/cuda/version.txt | sed "s/CUDA Version \(.*\)\..*$/\1/g")
else
    CUDA_VERSION=$(nvidia-smi -q | grep "CUDA Version" | sed "s/^CUDA Version.*:[^0-9]*\(.*\)\.\(.*\)\.*$/\1.\2/g")
fi
ENV_NAME=$OS_NAME-cuda-$CUDA_VERSION
KATAGO_VERSIONS="1.7.0"
GPU_NUM=$(($(nvidia-smi -q | grep "Attached GPUs" | cut -d':' -f2)))
echo "System Env: " $ENV_NAME
echo "GPU Info: " $GPU_NAME x $GPU_NUM
PACKAGE=ubuntu-cuda-$CUDA_VERSION
#if [ "$1" != "" ]
#then
#    PACKAGE=$1
#fi
if [[ "$PACKAGE" == "ubuntu-cuda-11"* ]]
then
    # use cuda 11.1 instead
    PACKAGE=ubuntu-cuda-11.1
fi
echo "USING PACKAGE: $PACKAGE"
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
        rm -rf $FILE_PATH
        mv $FILE_PATH.downloading $FILE_PATH
    fi
}
mkdir -p ./resources
echo "Downloading engines..."
for KATAGO_VERSION in $KATAGO_VERSIONS
do
    update_file ./resources/gomoku-$KATAGO_VERSION-$PACKAGE.zip $DATA_DOWNLOAD_URL/gomoku-$KATAGO_VERSION-$PACKAGE.zip
    if [ $? -ne 0 ]
    then
        echo "Failed to download the engines."
        exit -1
    fi
done

echo "Downloading weights..."
update_file ./resources/gomoku-weights.zip $DATA_DOWNLOAD_URL/gomoku-weights.zip
if [ $? -ne 0 ]
then
    echo "Failed to download the weights."
    exit -1
fi

echo "Downloading work..."
update_file ./resources/gomoku-work.zip $DATA_DOWNLOAD_URL/gomoku-work.zip 7306ece9dea85e619ec3c3484796af7c
if [ $? -ne 0 ]
then
    echo "Failed to download the work."
    exit -1
fi
echo "Download done"
echo "Installing..."

cd ./resources 
echo "Installing weights..."
rm -rf gomoku-weights && unzip gomoku-weights.zip >/dev/null 2>&1
echo "Installing engines..."
for KATAGO_VERSION in $KATAGO_VERSIONS
do
    rm -rf gomoku-$KATAGO_VERSION-$PACKAGE && unzip gomoku-$KATAGO_VERSION-$PACKAGE.zip  >/dev/null 2>&1
done
echo "Installing others..."
rm -rf gomoku-work && unzip gomoku-work.zip >/dev/null 2>&1
rm -rf configs && unzip configs.zip >/dev/null 2>&1
cd -
rm -rf work
mv ./resources/gomoku-work ./work
mkdir -p ./work/data
mv ./resources/gomoku-weights ./work/data/weights

mkdir -p ./work/data/configs
rm -rf ./work/data/configs/default_gtp.cfg
mv ./work/default.cfg ./work/data/configs/default_gtp.cfg
SEARCH_THREAD_NUM=$(($GPU_NUM * 48))
echo "Using Thread Num: $SEARCH_THREAD_NUM"

CMD="s/numSearchThreads = .*/numSearchThreads = $SEARCH_THREAD_NUM/g"
sed -i "$CMD" ./work/data/configs/default_gtp.cfg
CMD="s/^.*numNNServerThreadsPerModel.*$//g"
sed -i "$CMD" ./work/data/configs/default_gtp.cfg
CMD="s/^.*gpuToUseThread.*$//g"
sed -i "$CMD" ./work/data/configs/default_gtp.cfg

echo "numNNServerThreadsPerModel = $GPU_NUM" >> ./work/data/configs/default_gtp.cfg
i=0
while [ $i -lt $GPU_NUM ]
do
    echo "gpuToUseThread$i = $i" >> ./work/data/configs/default_gtp.cfg
    i=$(($i + 1))
done


mkdir -p ./work/data/bins
for KATAGO_VERSION in $KATAGO_VERSIONS
do
    mv ./resources/gomoku-$KATAGO_VERSION-$PACKAGE ./work/data/bins/gomoku-$KATAGO_VERSION
done

chmod +x ./work/*.sh
chmod +x ./work/ikatago-server

echo "Install Successfully, now you can run the ikatago-server"
