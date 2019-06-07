FROM gcr.io/com-seankhliao/imagemagick AS build-img

WORKDIR /workspace
RUN mkdir -p dst
RUN convert src/icon.tif -flatten \
        '(' +clone -resize 512x512 -quality 60 -write dst/icon-512.png +delete ')' \
        '(' +clone -resize 192x192 -quality 60 -write dst/icon-192.png +delete ')' \
        '(' +clone -resize 128x128 -quality 60 -write dst/icon-128.png +delete ')' \
        '(' +clone -resize 64x64 -quality 60 -write dst/icon-64.png +delete ')' \
        '(' +clone -resize 48x48 -quality 60 -write dst/icon-48.png +delete ')' \
        '(' +clone -resize 32x32 -quality 60 -write dst/icon-32.png +delete ')' \
        '(' +clone -resize 16x16 -quality 60 -write dst/icon-16.png +delete ')' \
        -resize 48x48 -quality 60 dst/favicon.ico
RUN convert -background none -density 1200 -resize 2100x1350 src/map.svg -write dst/map.webp dst/map.png


FROM gcr.io/com-seankhliao/parcel AS build-js

WORKDIR /workspace
RUN parcel build -d src/readss src/readss-src/index.js


FROM gcr.io/com-seankhliao/site-builder AS build-site

WORKDIR /workspace
RUN /bin/app
