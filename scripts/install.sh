#!/bin/bash
OS_NAME=$(cat /etc/os-release | grep "PRETTY_NAME" | sed 's/PRETTY_NAME="\(.*\)"/\1/g')

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
KATAGO_VERSIONS="1.8.1"
GPU_NUM=$(($(nvidia-smi -q | grep "Attached GPUs" | cut -d':' -f2)))
echo "System Env: " $ENV_NAME
echo "GPU Info: " $GPU_NAME x $GPU_NUM
PACKAGE=ubuntu-cuda-$CUDA_VERSION
#if [ "$1" != "" ]
#then
#    PACKAGE=$1
#fi
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
    update_file ./resources/katago-$KATAGO_VERSION-$PACKAGE.zip https://ikatago-resources.oss-cn-beijing.aliyuncs.com/all/katago-$KATAGO_VERSION-$PACKAGE.zip
    if [ $? -ne 0 ]
    then
        echo "Failed to download the engines."
        exit -1
    fi
done

echo "Downloading weights..."
update_file ./resources/weights.zip https://ikatago-resources.oss-cn-beijing.aliyuncs.com/all/weights.zip 8e633a8706d3a40781863720035102c9
if [ $? -ne 0 ]
then
    echo "Failed to download the weights."
    exit -1
fi
echo "Downloading configs..."
update_file ./resources/configs.zip https://ikatago-resources.oss-cn-beijing.aliyuncs.com/all/configs.zip 43cc69565e7cdf768047c70d0c1ee0f4
if [ $? -ne 0 ]
then
    echo "Failed to download the configs."
    exit -1
fi
echo "Downloading work..."
update_file ./resources/linux-work.zip https://ikatago-resources.oss-cn-beijing.aliyuncs.com/all/linux-work.zip bd0e3408ada19e666ef916daafd6dc5f
if [ $? -ne 0 ]
then
    echo "Failed to download the work."
    exit -1
fi
echo "Download done"
echo "Installing..."

cd ./resources 
echo "Installing weights..."
rm -rf weights && unzip weights.zip >/dev/null 2>&1
echo "Installing engines..."
for KATAGO_VERSION in $KATAGO_VERSIONS
do
    rm -rf katago-$KATAGO_VERSION-$PACKAGE && unzip katago-$KATAGO_VERSION-$PACKAGE.zip  >/dev/null 2>&1
done
echo "Installing others..."
rm -rf linux-work && unzip linux-work.zip >/dev/null 2>&1
rm -rf configs && unzip configs.zip >/dev/null 2>&1
cd -
rm -rf work
mv ./resources/linux-work ./work
mkdir -p ./work/data
mv ./resources/weights ./work/data/
mkdir -p ./work/data/configs
if [ $GPU_NUM -eq 1 ]
then
    cardsfolder=$GPU_NUM"card"
else
    cardsfolder=$GPU_NUM"cards"
fi

if [ -d "./resources/configs/$cardsfolder-$GPU_NAME" ]
then 
    cardsfolder="$cardsfolder-$GPU_NAME"
fi
cp ./resources/configs/$cardsfolder/* ./work/data/configs/
mkdir -p ./work/data/bins
for KATAGO_VERSION in $KATAGO_VERSIONS
do
    mv ./resources/katago-$KATAGO_VERSION-$PACKAGE ./work/data/bins/katago-$KATAGO_VERSION
done

echo "Install Successfully, now you can run the ikatago-server"
