make:
	go run main.go examples/slides.md

build:
	go build -o slides

install:
	go install
