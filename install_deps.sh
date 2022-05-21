IMAGEMAGICK_PROJECT=ImageMagick
IMAGEMAGICK_VERSION=7.1.0-31
IMAGEMAGICK_VERSION=$IMAGEMAGICK_VERSION

sudo apt-get install -y \
	wget build-essential \
	pkg-config \
	--no-install-recommends \
&& sudo apt-get -q -y install \
	libjpeg-dev \
	libpng-dev \
	libtiff-dev \
	libgif-dev \
	libx11-dev \
	fontconfig fontconfig-config libfontconfig1-dev \
	ghostscript gsfonts gsfonts-x11 \
	libfreetype6-dev \
	--no-install-recommends \
&& sudo rm -rf /var/lib/apt/lists/*

cd && \
	wget https://github.com/ImageMagick/${IMAGEMAGICK_PROJECT}/archive/refs/tags/${IMAGEMAGICK_VERSION}.tar.gz && \
	tar xvzf ${IMAGEMAGICK_VERSION}.tar.gz && \
	cd ImageMagick* && \
	./configure \
	    --without-magick-plus-plus \
	    --without-perl \
	    --disable-openmp \
	    --with-gvc=no \
	    --with-fontconfig=yes \
	    --with-freetype=yes \
	    --with-gslib \
	    --disable-docs && \
	make -j$(nproc) && sudo make install && \
	sudo ldconfig /usr/local/lib
