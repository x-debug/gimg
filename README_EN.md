![GImg](./resources/logo-192x192.png)

**Lightweight image processing service built by Golang**

GImg is a general-purpose image processing service that can scale images, make thumbnails, and more.

It provides a lot of object storage backends (localfs/memcached/redis/ssdb/fastdfs...), focusing on image processing.

A common pattern is to put it behind an HTTP proxy server, such as Nginx or Varnish.

### Setup
##### Source installation
```shell
git clone https://github.com/x-debug/gimg.git

cd gimg && make build

./gimg
```

##### Docker image installation
```shell
docker pull nbboy/gimg:v1.5

docker run -d -p 8888:8888 nbboy/gimg:v1.5
```

### Features
- Scale image
- Image thumbnail
- Image rotation 
- Image watermark
- Image cropping
- Image grayscale
- Image rounding
- Set image quality 
- Set image format
- Custom lua script

### Architecture design

### Some optimizations

##### Local storage design

##### Cache design

##### Performance Testing
