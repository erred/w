.PHONY: all
all: clean build

.PHONY: build
build:
	mkdir dist
	cp static/* dist
	cp src/*.html dist
	convert -background none -density 1200 -resize 2100x1350 \
		src/map.svg \
		-write dist/map.webp \
		dist/map.png
	convert src/icon.tif -flatten \
		\( +clone -resize 512x512 -quality 60 -write dist/icon-512.png +delete \) \
		\( +clone -resize 256x256 -quality 60 -write dist/icon-256.png +delete \) \
		\( +clone -resize 192x192 -quality 60 -write dist/icon-192.png +delete \) \
		\( +clone -resize 180x180 -quality 60 -write dist/icon-180.png +delete \) \
		\( +clone -resize 128x128 -quality 60 -write dist/icon-128.png +delete \) \
		\( +clone -resize 64x64 -quality 60 -write dist/icon-64.png +delete \) \
		\( +clone -resize 48x48 -quality 60 -write dist/icon-48.png +delete \) \
		\( +clone -resize 32x32 -quality 60 -write dist/icon-32.png +delete \) \
		\( +clone -resize 16x16 -quality 60 -write dist/icon-16.png +delete \) \
		-resize 16x16 -quality 60 dist/icon.ico

.PHONY: clean
clean:
	rm -rf ./dist

.PHONY: test
test:
