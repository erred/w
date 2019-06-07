FROM gcr.io/com-seankhliao/imagemagick AS build-img
WORKDIR /build
RUN convert /workspace/src/icon.tif -flatten \
        '(' +clone -resize 512x512 -quality 60 -write icon-512.png +delete ')' \
        '(' +clone -resize 192x192 -quality 60 -write icon-192.png +delete ')' \
        '(' +clone -resize 128x128 -quality 60 -write icon-128.png +delete ')' \
        '(' +clone -resize 64x64 -quality 60 -write icon-64.png +delete ')' \
        '(' +clone -resize 48x48 -quality 60 -write icon-48.png +delete ')' \
        '(' +clone -resize 32x32 -quality 60 -write icon-32.png +delete ')' \
        '(' +clone -resize 16x16 -quality 60 -write icon-16.png +delete ')' \
        -resize 48x48 -quality 60 favicon.ico
RUN convert -background none -density 1200 -resize 2100x1350 /workspace/src/map.svg -write map.webp map.png


FROM gcr.io/com-seankhliao/parcel AS build-js
WORKDIR /build
RUN parcel build -d /workspace/src/readss index.js


FROM gcr.io/com-seankhliao/site-builder AS build-site
WORKDIR /workspace/dst
COPY --from=build-img /build .
COPY --from=build-js /build /workspace/src/readss/
WORKDIR /workspace
RUN ["/bin/app"]
