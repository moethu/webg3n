all: build

run:
	docker run --rm -d -p 8000:8000 local/webg3n

build:
	docker build -t local/webg3n -f Dockerfile .
