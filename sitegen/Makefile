.PHONY: build
build:
	go build -mod=vendor

.PHONY: test
test:
	go test -mod=vendor

.PHONY: update
update:
	go get -u
	go mod tidy
	go mod vendor
