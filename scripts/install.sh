#!/bin/bash
OS_NAME=$(cat /etc/os-release | grep "PRETTY_NAME" | sed 's/PRETTY_NAME="\(.*\)"/\1/g')
if [[ "$OS_NAME" = "Ubuntu 18"* ]]
then
    OS_NAME=ubuntu-18
elif [[ "$OS_NAME" = "Ubuntu 16"* ]]
then
    OS_NAME=ubuntu-16
else
    echo "only ubuntu 16 or ubuntu 18 are supported"
    exit -1
fi

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
KATAGO_VERSION=1.6.1
GPU_NUM=$(($(nvidia-smi -q | grep "Attached GPUs" | cut -d':' -f2)))
echo "System Env: " $ENV_NAME
echo "GPU Info: " $GPU_NAME x $GPU_NUM
PACKAGE=$ENV_NAME
if [ "$GPU_NAME" == "A100" ] || [ "$GPU_NAME" == "3090" ]
then
    if [ "$CUDA_VERSION" != "11.1" ]
    then
        PACKAGE=$OS_NAME-cuda-"11.1-with-cuda"
    fi
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
        mv $FILE_PATH.downloading $FILE_PATH
    fi
}
mkdir -p ./resources
echo "Downloading engines..."
update_file ./resources/katago-$KATAGO_VERSION-$PACKAGE.zip https://ikatago-resources.oss-cn-beijing.aliyuncs.com/all/katago-$KATAGO_VERSION-$PACKAGE.zip
if [ $? -ne 0 ]
then
    echo "Failed to download the engines."
    exit -1
fi

echo "Downloading weights..."
update_file ./resources/weights.zip https://ikatago-resources.oss-cn-beijing.aliyuncs.com/all/weights.zip 5268d0ee93bbff1be25b472beb3c0022
if [ $? -ne 0 ]
then
    echo "Failed to download the weights."
    exit -1
fi
echo "Downloading configs..."
update_file ./resources/configs.zip https://ikatago-resources.oss-cn-beijing.aliyuncs.com/all/configs.zip 453070e588e1bcb5e23993b77264e8a2
if [ $? -ne 0 ]
then
    echo "Failed to download the configs."
    exit -1
fi
echo "Downloading work..."
update_file ./resources/linux-work.zip https://ikatago-resources.oss-cn-beijing.aliyuncs.com/all/linux-work.zip 6cce6794adf1d66c29878d0f6818e16d
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
rm -rf katago-$KATAGO_VERSION-$PACKAGE && unzip katago-$KATAGO_VERSION-$PACKAGE.zip  >/dev/null 2>&1
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
mv ./resources/katago-$KATAGO_VERSION-$PACKAGE ./work/data/bins/katago-$KATAGO_VERSION

echo "Install Successfully, now you can run the ikatago-server"
