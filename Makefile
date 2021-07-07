# Makefile for the project to keep the things simpler
export APP_DB_HOST=localhost
export APP_DB_USERNAME=postgres
export APP_DB_PASSWORD=password
export APP_DB_NAME=postgres

# -- COLORS
export STOP_COLOR="\e[0m"
# color for a main activity
export ACTIVITY="\e[1;34m➜\e[1;97m"
# color for a sub activity
export SUB_ACT="\n\e[0;33m"
export DONE="[\e[1;32mDONE\e[0m]"
export ERROR="[\e[1;31mERROR\e[0m]"

# start the application locally using:
# - docker container for postgres
# requirements for this test are: docker, golang sdk
run: rm-postgres-container
	@docker-compose up -d postgresql
	@go run *.go

# start a docker composition that has:
# - docker container for postgres
# - docker container for rest-api at the latest version
# requirements: docker
# Using this you can make an HTTP request on localhost:9090
run-docker-env: rm-postgres-container
	@docker-compose up -d --build

stop-docker-env:
	@docker-compose down -v
	@docker-compose rm -v

# test the application using:
# - docker container for postgres
# - go test
# requirements for this test are: docker, golang sdk
test-app:
	@clear
	@printf "\e[1;34m➜\e[1;97m starting tests using go and postgresql...\n\e[0m"
	@printf "\e[0;33mstarting postgresql as a docker container...\n\e[0m"
	@docker-compose up -d postgresql
	@sleep 1
	@export APP_CONFIG="${PWD}/config/server.yml"
	@printf "\e[0;33mstarting go test...\n\e[0m"
	@APP_CONFIG="${PWD}/config/server.yml" go test github.com/mas2020-golang/rest-api/test/...

# test the application using:
# - docker container for postgres (attached to a test network)
# - docker container run using the image built with the Dockerfile.test. This container is a golang runtime container.
# requirements: docker
test-docker-app: stop-docker
	@./test.sh

rm-postgres-container:
	@docker container rm -f postgres_test

clean:
	@printf "${PWD}\n"

