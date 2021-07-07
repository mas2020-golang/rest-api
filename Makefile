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

start:
	@go run *.go

test-app:
	@printf "\e[1;34m➜\e[1;97m starting tests...\n\e[0m"
	@printf "\e[0;33mstarting postgresql as a docker container...\n\e[0m"
	@docker-compose up -d postgresql
	@export APP_CONFIG="${PWD}/config/server.yml"
	@APP_CONFIG="${PWD}/config/server.yml" go test github.com/mas2020-golang/rest-api/test/...

test-docker-app:
	@./test.sh

clean:
	@printf "${PWD}\n"