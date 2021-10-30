#!/bin/bash
DATA_DOWNLOAD_URL="http://8.130.167.86:8080/api/ikatago/resources/download/all"

BACKEND_TYPE=$1
if [ "$BACKEND_TYPE" == "" ]
then
    BACKEND_TYPE="CUDA"
fi

THREAD_NUM=$2
if [ "$THREAD_NUM" == "" ]
then
    THREAD_NUM="AUTO"
fi

OS_NAME=$(cat /etc/os-release | grep "^NAME=" | sed 's/NAME="\(.*\)"/\1/g')
OS_VERSION=$(cat /etc/os-release | grep "^VERSION_ID=" | sed 's/VERSION_ID="\(.*\)"/\1/g' | tr -d [.])
if [ "$OS_NAME" != "Ubuntu" ]
then
    echo "Failed to installed. Only Ubuntu 18.04 or Ubuntu 20.04 are supported. Current is " + $OS_NAME + " " + $OS_VERSION
    exit -1
fi
OS_NAME="ubuntu"

if [ "$OS_VERSION" != "1804" -a "$OS_VERSION" != "2004" ]
then
    echo "Failed to installed. Only Ubuntu 18.04 or Ubuntu 20.04 are supported. Current is " + $OS_NAME + " " + $OS_VERSION
    exit -1
fi


if [ -f /usr/local/cuda/version.txt  ]
then
    CUDA_VERSION=$(cat /usr/local/cuda/version.txt | sed "s/CUDA Version \(.*\)\..*$/\1/g")
else
    CUDA_VERSION=$(nvidia-smi -q | grep "CUDA Version" | sed "s/^CUDA Version.*:[^0-9]*\(.*\)\.\(.*\)\.*$/\1.\2/g")
fi
if [[ "$CUDA_VERSION" != "11."* ]]
then
   echo "Only Cuda 11.x are supported. Current is: " + $CUDA_VERSION
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
elif [[ "$GPU_NAME" == *"A6000"* ]]
then
    GPU_NAME="A6000"
fi

ENV_NAME=$OS_NAME-$OS_VERSION-cuda-$CUDA_VERSION
KATAGO_VERSIONS="1.10"
GPU_NUM=$(($(nvidia-smi -q | grep "Attached GPUs" | cut -d':' -f2)))
echo "System Env: " $ENV_NAME
echo "GPU Info: " $GPU_NAME x $GPU_NUM
PACKAGE=ubuntu-20-cuda-11.4-$BACKEND_TYPE
#PACKAGE=ubuntu-cuda-11.1

echo "USING PACKAGE: $PACKAGE"
update_file() {
    FILE_PATH=$1
    FILE_URL=$2
    TARGET_MD5=$3
    SHOULD_UPDATE=0
    #echo $2
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
    update_file ./resources/katago-$KATAGO_VERSION-$PACKAGE.zip $DATA_DOWNLOAD_URL/katago-$KATAGO_VERSION-$PACKAGE.zip
    if [ $? -ne 0 ]
    then
        echo "Failed to download the engines.: katago-$KATAGO_VERSION-$PACKAGE.zip"
        exit -1
    fi
done

echo "Downloading weights..."
update_file ./resources/weights.zip $DATA_DOWNLOAD_URL/weights.zip 6154da3f872cbbdce4227772e998a6c0
if [ $? -ne 0 ]
then
    echo "Failed to download the weights."
    exit -1
fi
echo "Downloading configs..."
update_file ./resources/configs.zip $DATA_DOWNLOAD_URL/configs.zip 68f1ea2261598609cfd9cfec270c8755
if [ $? -ne 0 ]
then
    echo "Failed to download the configs."
    exit -1
fi
echo "Downloading work..."
update_file ./resources/linux-work.zip $DATA_DOWNLOAD_URL/linux-work.zip dc4a0274a2bc0fbcbd3ca2ee1fddf07c
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

if [ "$BACKEND_TYPE" == "TENSORRT" ]
then
    #install tensorrt
    apt-key adv --fetch-keys https://ikatago-resources.oss-cn-beijing.aliyuncs.com/all/7fa2af80.pub
    echo "deb http://developer.download.nvidia.com/compute/cuda/repos/ubuntu${OS_VERSION}/x86_64 /" > /etc/apt/sources.list.d/cuda.list
    echo "deb http://developer.download.nvidia.com/compute/machine-learning/repos/ubuntu${OS_VERSION}/x86_64 /" > /etc/apt/sources.list.d/cuda_learn.list
    apt update
    apt-get install -y libnvinfer8=8.2.0-1+cuda11.4
    # modify config
    sed -i "s/^cudaDeviceToUseThread\(.*\)$/trtDeviceToUseThread\1/g" ./work/data/configs/*.cfg
fi

if [ "$THREAD_NUM" != "AUTO" ]
then
    sed -i "s/^numSearchThreads = .*$/numSearchThreads = $THREAD_NUM/g" ./work/data/configs/*.cfg
fi

echo "Install Successfully, now you can run the ikatago-server"
