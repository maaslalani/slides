make:
	go run main.go examples/slides.md

test:
	go test ./... -short

build:
	go build -o slides

install:
	go install
