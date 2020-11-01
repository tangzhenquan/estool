GOPATH:=$(shell go env GOPATH)

PROG_NAME:=esImport
.PHONY: build
build:
	go build -o $(PROG_NAME) main.go
.PHONY: test
test:
	go test -v ./... -cover

.PHONY: install
install: build
	mv $(PROG_NAME) $(DESTDIR)
	cp esImport.yaml $(DESTDIR)