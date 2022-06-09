<img align="right" width="150px" src="./resources/logo-192x192.png">

# Gimg
[English](README_EN.md) | 简体中文

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Golang构建的轻量级图片处理服务**

Gimg是通用的图片处理服务,能够缩放图片,制作缩略图等。

它提供非常多的对象存储后端(localfs/fastdfs...),对热点图片进行缓存(memory/memcached/redis/ssdb),专注做图片的处理。

您可以把它放到HTTP代理服务器之后,比如Nginx或者Varnish,也可以指定图片获取的路径,它会自动去抓取图片。

**项目目前没有经过大量测试,切勿使用在生产环境中**

### 安装
##### 源码安装
```shell
git clone https://github.com/x-debug/gimg.git

cd gimg && make build

./gimg
```

##### Docker镜像安装
```shell
docker pull nbboy/gimg:v1.5.2

docker run -d -p 8888:8888 nbboy/gimg:v1.5.2
```

##### 上传测试
安装完成后,打开http://YouDomain:8888/demo 可以上传测试的文件

### 功能
- 图片缩放
```
http://YouDomain:8888/图片Hash?op=resize&w=100&h=100
```
- 图片缩略图
```
http://YouDomain:8888/图片Hash?op=thumbnail&w=100&h=100
```
- 图片旋转
```
http://YouDomain:8888/图片Hash?op=rotate&deg=30
```
- 图片水印
- 图片裁剪
```
http://YouDomain:8888/图片Hash?op=crop&x=30&y=60&w=300&h=300
```
- 图片灰度化
```
http://YouDomain:8888/图片Hash?op=grayscale
```
- 图片圆角处理
```
http://YouDomain:8888/图片Hash?op=round&rx=30&ry=30
```
- 图片压缩处理
```
http://YouDomain:8888/图片Hash?op=quality&q=30
```
- 图片格式处理
```
http://YouDomain:8888/图片Hash?op=format&f=png
```
- 自定义LUA脚本处理
```
http://YouDomain:8888/图片Hash?op=lua&f=demo
```
- 图片回源处理
```
http://YouDomain:8888/?remote=https://alifei05.cfp.cn/creative/vcg/veer/1600water/veer-140775274.jpg&op=rotate&deg=30
```
### 架构设计
进行中

### 存储目录设计
根据MD5的前六位进行哈希，1-3位转换为十六进制数后除以4，范围正好落在1024以内，以这个数作为第一级子目录；4-6位同样处理，作为第二级子目录；二级子目录下是以MD5命名的文件夹，每个MD5文件夹内存储图片的原图和其他根据需要存储的版本。

### 缓存设计
缓存主要对处理过后的图片进行保存,以加速返回的速度。缓存主要分两级,分为第一级缓存和第二级缓存,第一级缓存可以配置,目前支持Memcache/Memory缓存,当然也可以关闭,第二级缓存为磁盘缓存,只有在第一级缓存没命中才会去取第二级缓存。如果第一二级缓存都没命中,则最后才去处理原图片。


### 性能测试
进行中

**如果觉得项目有帮助，请给一个免费的小红心⭐**
