# This is a veeery simplistic makefile

.PHONY: build

default:
	go run .

build: 
	mkdir -p ./build && go build -o ./build .

run:
	go run .

clean:
	rm ./build/* ; rm -rf ./.jorge

install:
	cp ./build/jorge /usr/bin/jorge

uninstall:
	rm /usr/bin/jorge
