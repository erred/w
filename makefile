.PHONY: build deploy serve

build:
	inkscape -z -D -y 0 -e map.png -w 2100 -h 1350 map.svg

deploy:
	firebase deploy

