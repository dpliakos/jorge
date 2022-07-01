.PHONY: build

default:
	go run .

build: 
	mkdir -p ./build && go build -o ./build .

run:
	go run .

clean:
	rm ./build/* ; rm -rf ./.jorge
