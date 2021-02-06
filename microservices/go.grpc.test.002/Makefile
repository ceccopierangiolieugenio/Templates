.PHONY: build run

build:
	docker-compose build

run:
	docker-compose rm -f
	docker-compose up --force-recreate --remove-orphans --abort-on-container-exit

unit_test:
	docker-compose -f test/docker-compose.test.backend.yml --project-directory .  up  --abort-on-container-exit --remove-orphans --build
	docker-compose -f test/docker-compose.test.server.yml  --project-directory .  up  --abort-on-container-exit --remove-orphans --build

postman_test:
	docker-compose -f test/docker-compose.postman.test.yml --project-directory .  up  --abort-on-container-exit --remove-orphans --build

test: unit_test postman_test