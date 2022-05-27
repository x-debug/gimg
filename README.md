![Gimg](./resources/logo-192x192.png)

**Golang构建的轻量级图片处理服务**

Gimg是通用的图片处理服务,能够缩放图片,制作缩略图等。

它提供非常多的对象存储后端(localfs/fastdfs...),对热点图片进行缓存(memory/memcached/redis/ssdb),专注做图片的处理。

您可以把它放到HTTP代理服务器之后,比如Nginx或者Varnish,也可以指定图片获取的路径,它会自动去抓取图片。

### 安装
##### 源码安装
```shell
git clone https://github.com/x-debug/gimg.git

cd gimg && make build

./gimg
```

##### Docker镜像安装
```shell
docker pull nbboy/gimg:v1.5

docker run -d -p 8888:8888 nbboy/gimg:v1.5
```

### 功能
- 图片缩放
- 图片缩略图
- 图片旋转
- 图片水印
- 图片裁剪
- 图片灰度化
- 图片圆角处理
- 图片压缩处理
- 图片格式处理
- 自定义LUA脚本处理

### 架构设计

### 一些优化

##### 本地目录存储设计

##### 缓存设计

##### 图片处理设计

##### 性能测试
