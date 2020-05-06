.PHONY: *

OUT := public

all: build

clean:
	rm -rf ${OUT}

build:
	webrender
	# @docker run --rm -it -v $$(pwd):/workspace -w /workspace seankhliao/webstyle

img: img-icon img-map

img-icon:
	@mkdir -p ${OUT} || true
	@docker run --rm -it -v $$(pwd):/workspace -w /workspace seankhliao/ci-imagemagick \
		img/icon.tif -flatten \
		"(" +clone -resize 512x512 -quality 60 -write ${OUT}/icon-512.png +delete ")" \
		"(" +clone -resize 192x192 -quality 60 -write ${OUT}/icon-192.png +delete ")" \
		"(" +clone -resize 128x128 -quality 60 -write ${OUT}/icon-128.png +delete ")" \
		"(" +clone -resize 64x64 -quality 60 -write ${OUT}/icon-64.png +delete ")" \
		"(" +clone -resize 48x48 -quality 60 -write ${OUT}/icon-48.png +delete ")" \
		"(" +clone -resize 32x32 -quality 60 -write ${OUT}/icon-32.png +delete ")" \
		"(" +clone -resize 16x16 -quality 60 -write ${OUT}/icon-16.png +delete ")" \
		-resize 32x32 ${OUT}/favicon.ico

img-map:
	@mkdir -p ${OUT} || true
	@docker run --rm -it -v $$(pwd):/workspace -w /workspace seankhliao/ci-imagemagick \
		 -background none -density 1200 -resize 1920x1080 img/map.svg \
		-write ${OUT}/map.png \
		-write ${OUT}/map.webp \
		${OUT}/map.jpg


update:
	docker pull seankhliao/ci-imagemagick:latest
	docker pull seankhliao/webstyle:latest
