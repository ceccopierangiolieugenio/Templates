.PHONY: build run

build:
	docker-compose build

run:
	docker-compose up --scale backend=3 --scale web_server=3 