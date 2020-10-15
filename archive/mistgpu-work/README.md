# mist上的安装办法 (cuda-10.1环境)
1). 在mistgpu的服务器上，安装必要的库, 运行:

```
sudo apt update && sudo apt-get install -y libzip-dev zlib1g-dev libboost-filesystem-dev
```

2). 
```
unzip mistgpu-work.zip 
```

3). 在mistgpu的网站上，解压这个包，然后进入`mistgpu-work`目录
4). 在mistgpu上找到你的公开端口（不是ssh端口），可以在"帮助说明“最下面看到你的服务器的公开端口.
5). 修改`config/conf.yaml`里的这两个地方:
```
server:
  listen: 0.0.0.0:55504 # 55504 -> 改成你自己的公开端口
```
```
  direct: 
    type: direct
    host: gpu48.mistgpu.com # gpu48.mistgpu.com -> 改成你的mistgpu的主机名
    port: 55504 # 55504 -> 改成你自己的公开端口
```

3). 运行
```
./run.sh 你的用户名 你的密码
```