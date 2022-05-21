![GImg](./resources/logo-192x192.png)

**Golang构建的轻量级图片处理服务**

GImg是通用的图片处理服务，能够进行图片缩放，制作缩略图，等等。

它提供非常多的对象存储后端（localfs/memcached/redis/ssdb/fastdfs...），专注做图片的处理。

常用的模式是把它放到HTTP代理服务器之后，比如Nginx或者Varnish。

### 安装
##### 源码安装
```shell
git clone https://github.com/x-debug/gimg.git

cd gimg && make build

./gimg
```

##### Docker镜像安装
```shell
docker pull nbboy/gimg:v1.4

docker run -d -p 8888:8888 nbboy/gimg:v1.4
```

### 功能
- 图片缩放
- 图片缩略图
- 图片旋转
- 图片水印
- 图片裁剪
- 图片灰度化
- 图片圆角处理
- LUA图片处理定制

### 架构设计

### 一些优化

##### 本地目录存储设计

##### 缓存设计

##### 图片处理设计

##### 性能测试
