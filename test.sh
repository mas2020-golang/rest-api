# -- COLORS
export STOP_COLOR="\e[0m"
# color for a main activity
export ACTIVITY="\e[1;34m➜\e[1;97m"
# color for a sub activity
export SUB_ACT="\n\e[93m➜"
export DONE="[\e[1;32mDONE\e[0m]"
export ERROR="[\e[1;31mERROR\e[0m]"

clear
printf "${ACTIVITY} test rest-api using docker containers...${STOP_COLOR}\n"
printf "${SUB_ACT} creating the rest-api-test network...${STOP_COLOR}\n"
docker network create --attachable rest-api-test

printf "${SUB_ACT} starting postgresql container...${STOP_COLOR}\n"
docker run --rm -d -p 5432:5432 -e POSTGRES_PASSWORD=password \
-v "${PWD}"/scripts/db:/docker-entrypoint-initdb.d/ \
--network rest-api-test \
--name postgres_test postgres:13.2-alpine

printf "${SUB_ACT} building the latest docker image for the rest-api test...${STOP_COLOR}\n"
docker build -t appway/rest-api-test:latest -f Dockerfile.test .

printf "${SUB_ACT} executing the test...${STOP_COLOR}\n"
docker run --rm -it \
-e APP_DB_HOST=postgres_test \
-e APP_DB_USERNAME=postgres \
-e APP_DB_PASSWORD=password \
-e APP_DB_NAME=postgres \
-e APP_CONFIG=/usr/local/rest-api/config/server.yml \
--network rest-api-test \
--name rest-api-test appway/rest-api-test:latest

