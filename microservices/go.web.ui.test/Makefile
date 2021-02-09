.PHONY: build run

build:
	docker-compose build

run:
	docker-compose up --remove-orphans

unit_test:
	docker-compose -f test/docker-compose.test.server.yml  --project-directory .  up  --abort-on-container-exit --remove-orphans --build

postman_test:
	docker-compose -f test/docker-compose.postman.test.yml --project-directory .  up  --abort-on-container-exit --remove-orphans --build

robot_test:
	docker-compose -f test/docker-compose.robotframework.test.yml --project-directory .  up  --abort-on-container-exit --remove-orphans --build

test: unit_test postman_test robot_test