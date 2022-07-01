default:
	go run .

build: 
	mkdir -p ./build && go build -o ./build .

run:
	go run .
