.DEFAULT_GOAL := all
.PHONY: all build clean deploy setup

all: clean build

build:
	gotx -src src -out public
	rmdir public/partial

clean:
	rm -rf public/*

deploy:
	firebase deploy

setup:
	go get -u github.com/seankhliao/gotx
	firebase login


