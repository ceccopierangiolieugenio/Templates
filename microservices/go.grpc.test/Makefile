.PHONY: build run

build:
	docker-compose build

run:
	docker-compose up --remove-orphans --scale backend=3 --scale web_server=3 

unit_test:
	docker-compose -f test/docker-compose.test.backend.yml --project-directory .  up  --abort-on-container-exit --remove-orphans --build
	docker-compose -f test/docker-compose.test.server.yml  --project-directory .  up  --abort-on-container-exit --remove-orphans --build

functional_test:
	docker-compose -f test/docker-compose.functional.test.yml --project-directory .  up  --abort-on-container-exit --remove-orphans --build

test: unit_test functional_test