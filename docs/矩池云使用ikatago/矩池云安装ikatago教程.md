https://github.com/kinfkong/ikatago-resources/tree/master/dockerfiles


![image](https://user-images.githubusercontent.com/16535685/132154463-e22e1a1a-6023-46bb-8ad3-95b3100879cf.png)


从作者的库中可以看到，该程序支持cuda9.2、cuda10、cuda10.1、cuda10.2、cuda11.1等镜像,矩池云上的镜像基本上都可以满足他的要求，可以任意选用。


## 案例：用的cuda10.2的镜像

![image](https://user-images.githubusercontent.com/16535685/132154492-a964845d-3dc0-4f04-8a0b-1fff5e79a88b.png)
![image](https://user-images.githubusercontent.com/16535685/132154509-058b5578-b67b-4111-b9d7-49579d756827.png)
![image](https://user-images.githubusercontent.com/16535685/132154526-31e6b4b6-66c9-4d5a-829a-82ed35ec49f5.png)

## 利用脚本安装

```bash
cd ~; /bin/bash -c "$(curl -fsSL https://ikatago-resources.oss-cn-beijing.aliyuncs.com/all/install.sh)"
```

如果报错“sh: curl: command not found”是没有curl，先安装一下。

```bash
apt-get update
apt install curl
```


![image](https://user-images.githubusercontent.com/16535685/132154653-46adf74b-b21c-4edb-aae1-43faa3fd2ebe.png)


安装后，文件路径在root目录下的work文件夹内，文件有如下


![image](https://user-images.githubusercontent.com/16535685/132154669-7a882811-453c-4b94-83b9-34879bebc5e6.png)


## 运行


运行命令:

```bash
cd ~/work; ./run.sh 你的用户名 你的密码
```

 **建议使用 挂后台运行命令:**

```bash
cd ~/work; nohup ./run.sh 你的用户名 你的密码 &
```

密码建议使用复杂一些的密码，可以用生成工具生成，比如lastpass的password generator。

https://www.lastpass.com/password-generator

案例如下：

```bash
cd ~/work; nohup ./run.sh matpool sNoeoLSVDVrZ &
```


![image](https://user-images.githubusercontent.com/16535685/132154693-02e5bca9-8c21-4691-8f83-a7675356715c.png)
![image](https://user-images.githubusercontent.com/16535685/132154704-59bbc34f-6ebb-44b7-b404-aad3ec9922c2.png)

## 下载ikatago-client

https://github.com/kinfkong/ikatago-client

https://github.com/kinfkong/ikatago-client/releases/tag/1.3.3

直接下载

https://github.com/kinfkong/ikatago-client/releases/download/1.3.3/ikatago-1.3.3-win64.zip


![image](https://user-images.githubusercontent.com/16535685/132154721-498e3ae1-0a0d-4d53-adf2-28150ba9b33a.png)

## 下载Sabaki

https://github.com/SabakiHQ/Sabaki

https://github.com/SabakiHQ/Sabaki/releases/tag/v0.51.1


![image](https://user-images.githubusercontent.com/16535685/132154746-70c75176-1aa5-4414-82f0-7d60c98c5d66.png)


**portable便携版**

https://github.com/SabakiHQ/Sabaki/releases/download/v0.51.1/sabaki-v0.51.1-win-x64-portable.exe

**安装版**

https://github.com/SabakiHQ/Sabaki/releases/download/v0.51.1/sabaki-v0.51.1-win-x64-setup.exe


## Sabaki配置

![image](https://user-images.githubusercontent.com/16535685/132154760-f995f9b3-6b13-4c60-b597-b0fb6f964740.png)

在菜单栏点击 Engines - Show Engines Sidebar 显示侧边引擎栏。


![image](https://user-images.githubusercontent.com/16535685/132154773-d8a216ad-f3c0-4352-98ae-d67d74d1fa31.png)

引擎栏点击 Attach Engine... 按钮，选择 Manage Engines...。


![image](https://user-images.githubusercontent.com/16535685/132154789-1b65e000-214f-4f31-812e-06c461d4f745.png)


在引擎菜单中分别填写 4 行引擎信息。

引擎名称：自定义填写。

这里我写的是

```bash
matpool
```

路径：ikatago 客户端路径，可点击前方文件夹图标通过浏览选择。

```bash
D:\ikatago-1.3.3-win64\ikatago.exe
```

参数：ikatago 客户端参数，用户密码替换为服务端启动时的用户名和密码参数。

```bash
--platform all --username USER_NAME --password USER_PASSWORD
```

此次我的是

```bash
--platform all --username matpool --password sNoeoLSVDVrZ
```

初始命令：可定义一些命令参数，如定义 10 秒下一次棋。

```bash
time_settings 0 10 1
```

填写完成后点击 Close。



![image](https://user-images.githubusercontent.com/16535685/132154816-b7a99e85-6924-46a3-b7b3-0f7b669c312c.png)


引擎栏点击 Attach Engine... 按钮，选择刚创建的引擎点击他。



![image](https://user-images.githubusercontent.com/16535685/132154833-e7b1c062-2128-4215-a7c8-c52ccabbd15a.png)


点击 Start Engine vs. Engine Game 开始机机对弈，每过 10 秒机器会走出一步。再次点击该按钮可以停止对弈。



![image](https://user-images.githubusercontent.com/16535685/132154850-23771e5d-9c7f-4645-be31-894f911a60d7.png)



## 保存环境下次使用


![image](https://user-images.githubusercontent.com/16535685/132154868-efd22116-53c2-4b74-82e4-9ca8ec324498.png)
![image](https://user-images.githubusercontent.com/16535685/132154876-d08ad51d-297b-497d-a6e0-96b72d37967d.png)


这样下次可以直接使用，不用再配置环境了。


详情可以查看[如何使用矩池云的保存环境功能](https://matpool.com/supports/snapshot)。



## 参考文章

[在百度aistudio上跑katago (v100), 然后可以用Sabaki, Lizzie等进行远程连接。](https://aistudio.baidu.com/aistudio/projectdetail/1208079?channelType=0&channel=0)
[ikatago-server](https://github.com/kinfkong/ikatago-server/blob/main/docs/%E6%9E%81%E9%93%BE%E4%BA%91%E4%BD%BF%E7%94%A8ikatago/ikatago.ipynb)
