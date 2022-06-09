<img align="right" width="150px" src="./resources/logo-192x192.png">

# Gimg
[中文](README.md) | English 

**Lightweight image processing service built by Golang**

Gimg is a general-purpose image processing service that can scale images, make thumbnails, and more.

It provides a lot of object storage backends (localfs/fastdfs...), and cache hot images(memory/memcached/redis/ssdb), focusing on image processing.

We can put it behind an HTTP proxy server, such as Nginx or Varnish. You can specify the image path and it will grab it.

**This is just a toy project, don't use it in a production environment**

### Setup
##### Source installation
```shell
git clone https://github.com/x-debug/gimg.git

cd gimg && make build

./gimg
```

##### Docker image installation
```shell
docker pull nbboy/gimg:v1.5.2

docker run -d -p 8888:8888 nbboy/gimg:v1.5.2
```

##### Upload and testing
After the installation is complete, open http://YourDomain:8888/demo to upload the test file

### Features
- Scale image
```
http://YouDomain:8888/图片Hash?op=resize&w=100&h=100
```
- Image thumbnail
```
http://YouDomain:8888/图片Hash?op=thumbnail&w=100&h=100
```
- Image rotation 
```
http://YouDomain:8888/图片Hash?op=rotate&deg=30
```
- Image watermark
- Image cropping
```
http://YouDomain:8888/图片Hash?op=crop&x=30&y=60&w=300&h=300
```
- Image grayscale
```
http://YouDomain:8888/图片Hash?op=grayscale
```
- Image rounding
```
http://YouDomain:8888/图片Hash?op=round&rx=30&ry=30
```
- Set image quality 
```
http://YouDomain:8888/图片Hash?op=quality&q=30
```
- Set image format
```
http://YouDomain:8888/图片Hash?op=format&f=png
```
- Custom lua script
```
http://YouDomain:8888/图片Hash?op=lua&f=demo
```
- Process image with remote url
```
http://YouDomain:8888/?remote=https://alifei05.cfp.cn/creative/vcg/veer/1600water/veer-140775274.jpg&op=rotate&deg=30
```

### Architecture design
Doing

### Storage design
Doing

### Cache design
Doing

### Performance Testing
Doing

**If you like or are using this project start your solution, please give it a star⭐. Thanks!**
