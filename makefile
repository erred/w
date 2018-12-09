.PHONY: build deploy serve

build:
	yarn run polymer build

deploy:
	firebase deploy

serve:
	yarn run polymer serve --hostname 0.0.0.0
