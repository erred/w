.PHONY: all
all: clean build

.PHONY: build
build:
	mkdir dist
	cp src/* dist
	-cp images/*.png dist
	-cp images/*.ico dist
	-cp images/*.webp dist
	convert -background none -density 1200 -resize 2100x1350 \
		images/map.svg \
		-write dist/map.webp \
		dist/map.png

.PHONY: clean
clean:
	rm -rf ./dist

.PHONY: test
test:

.PHONY: icon
icon:
	convert images/icon.tif -flatten \
		\( +clone -resize 512x512 -quality 60 -write images/icon-512.png +delete \) \
		\( +clone -resize 192x192 -quality 60 -write images/icon-192.png +delete \) \
		\( +clone -resize 128x128 -quality 60 -write images/icon-128.png +delete \) \
		\( +clone -resize 64x64 -quality 60 -write images/icon-64.png +delete \) \
		\( +clone -resize 48x48 -quality 60 -write images/icon-48.png +delete \) \
		\( +clone -resize 32x32 -quality 60 -write images/icon-32.png +delete \) \
		\( +clone -resize 16x16 -quality 60 -write images/icon-16.png +delete \) \
		-resize 16x16 -quality 60 images/icon.ico
