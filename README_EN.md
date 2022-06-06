![Gimg](./resources/logo-192x192.png)

**Lightweight image processing service built by Golang**

Gimg is a general-purpose image processing service that can scale images, make thumbnails, and more.

It provides a lot of object storage backends (localfs/fastdfs...), and cache hot images(memory/memcached/redis/ssdb), focusing on image processing.

We can put it behind an HTTP proxy server, such as Nginx or Varnish. You can specify the image path and it will grab it.

### Setup
##### Source installation
```shell
git clone https://github.com/x-debug/gimg.git

cd gimg && make build

./gimg
```

##### Docker image installation
```shell
docker pull nbboy/gimg:v1.5.1

docker run -d -p 8888:8888 nbboy/gimg:v1.5.1
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
- Process image with remote url

### Architecture design
Doing

### Local storage design
Doing

### Cache design
Doing

### Performance Testing
Doing
