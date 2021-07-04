clear
printf "\e[93m➜ creating the rest-api-test network...\e[0m\n"
docker network create --attachable rest-api-test

printf "\n\e[93m➜ starting postgresql container...\e[0m\n"
docker run --rm -d -p 5432:5432 -e POSTGRES_PASSWORD=password \
-v "${PWD}"/scripts/db:/docker-entrypoint-initdb.d/ \
--network rest-api-test \
--name postgres_test postgres

printf "\n\e[93m➜ building the latest docker image for the rest-api test...\e[0m\n"
docker build -t appway/rest-api-test:latest -f Dockerfile.t .

printf "\n\e[93m➜ executing the test...\e[0m\n"
docker run --rm -it \
-e APP_DB_HOST=postgres_test \
-e APP_DB_USERNAME=postgres \
-e APP_DB_PASSWORD=password \
-e APP_DB_NAME=postgres \
--network rest-api-test \
--name rest-api-test appway/rest-api-test:latest

