# rest-api <!-- omit in toc -->

Example project to build a rest API server with Go. The scope of this application is to use the smallest numbers of
extra packages to realize a full REST api server.

## Table of Content <!-- omit in toc -->

- [Resources](#resources)
- [Dependencies](#dependencies)
- [Structure of the application](#structure-of-the-application)
- [Start the application](#start-the-application)
- [Test the application](#test-the-application)
- [Curl examples for the 'products' handler](#curl-examples-for-the-products-handler)
- [Run in a Docker container](#run-in-a-docker-container)

## Resources

This is the list of the reference to the resource that we are going to use in the application:

- **Gorilla Mux**, official documentation can be found [here](https://www.gorillatoolkit.org/). The Gorilla world on
  github is [here](https://github.com/gorilla).
- **Validation of JSON** object is done using the package ***validator***. The repository can be
  found [here](https://github.com/go-playground/validator).
- To connect to PostgreSQL I used the [pgx](https://pkg.go.dev/github.com/jackc/pgx) package

For the principles to follow when creating a REST API application take a look at
this [article](https://docs.microsoft.com/en-us/azure/architecture/best-practices/api-design).

## Dependencies

Before to start ensure to have the correct packages:

```shell
require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgx/v4 v4.11.0
)
```

## Structure of the application

The application has these folders:

- handlers: contains each http handler

## Start the application (for dev and test envs)

To start the application, first set the correct environment variable. These variables are needed for the DB connection:

```shell
export APP_DB_HOST=localhost \
export APP_DB_USERNAME=postgres \
export APP_DB_PASSWORD=password \
export APP_DB_NAME=postgres
```

then start a docker container to host postgres. The database will create all the needed tables simply executing the init
ddl script that we attach as a volume:

```shell
docker run --rm -d -p 5432:5432 -e POSTGRES_PASSWORD=password \
-v ${PWD}/scripts/db:/docker-entrypoint-initdb.d/ --name postgres_test postgres:13.2-alpine
```

finally execute the application typing:

```shell
go run *.go
```

to start having a file log you can use the redirection of the standard output:

```shell
go run *.go > rest-api.log
```

You can execute the Curl example calls to test the application.

### Other env variables

Important env variables, other than the DB as seen above, are:

- ***APP_CONFIG*** [optional]: represents the config file for the application. In case you do not pass the default value
  is `config/server.yml`.
- ***APP_JWTPWD*** [optional]: represents the pwd to use when signing the token. If the variable is not given the pwd
is created randomly.
  
### Test the application

To test, first add the environment variables, then execute:

```shell
export APP_CONFIG=/Users/andrea/development/go/test-projects/rest-api/config/server.yml
go test github.com/mas2020-golang/rest-api/test/...
```

to have more information add the `-v` flag.

## Curl examples for the 'products' handler

- **LOGIN** to the application:
    - *root* pwd is `my-root-pwd`
    - *andrea* pwd is `my-andrea-pwd`

```shell
curl -v -s -X POST http://localhost:9090/login \
-H "Content-Type: multipart/form-data" \
-d '
{
    "username": "andrea",
    "password": "my-andrea-pwd"
} 
' | jq
#sed 's+\([a-zA-Z0-9]*\.[a-zA-Z0-9]*\.[a-zA-Z0-9]*\).*+\1+'
```

- **GET** all the products

```shell
export TOKEN=
curl -v -s  http://localhost:9090/products \
-H "Authorization: Bearer ${token}" | jq
```

- **GET** the single product

```shell
curl -v -s  http://localhost:9090/products/1 \
-H "Authorization: Bearer ${token}" | jq
```

- **CREATE** a new product

```shell
curl -s -X POST http://localhost:9090/products \
-H "Authorization: Bearer ${token}" \
-d '
{
    "name": "Espresso 2",
    "description": "Short and strong coffee",
    "price": 2.50,
    "sku": "dfadds-das-fdsa"
}' | jq
```

- **UPDATE** an existing product

```shell
curl -s -i -X PUT http://localhost:9090/products/1 \
-H "Authorization: Bearer {token}" \
--models-binary @- << EOF
{
    "name": "Espresso 900",
    "description": "More than a coffee",
    "price": 2.99,
    "sku": "df-d-fdsa"
}
EOF
```

## Run as a Docker container

To execute the application as a Docker container we need to use a composition. This is intended for **_dev_** and **_
test_** only environments. By this way we benefit from having **_postgres_** up and running somewhere. Our composition
is made by:

- docker container for **postgres 13.2**
- docker container for our **rest-api server**

To run the composition you need to have `docker-compose` installed on the machine. To start the solution you can type:

```shell
docker-compose up --build
```

`--build` option is to ensure you have the latest version of the docker image for the rest api server.

### Executing test using Docker containers only

To execute test using only containers (for example in order to pass some continuous integration build step)
you have to create a network first:

```shell
docker network create --attachable rest-api-test
```

then you start your containers and attach to the right network (by this way the containers can see each other using the
name). First start the postgres database as:

```shell
docker run --rm -d -p 5432:5432 -e POSTGRES_PASSWORD=password \
-v ${PWD}/scripts/db:/docker-entrypoint-initdb.d/ \
--network rest-api-test \
--name postgres_test postgres
```

build the docker image for testing:

```shell
docker build -t appway/rest-api-test:latest -f Dockerfile.t .
```

then run the test as:

```shell
docker run --rm -it \
-e APP_DB_HOST=postgres_test \
-e APP_DB_USERNAME=postgres \
-e APP_DB_PASSWORD=password \
-e APP_DB_NAME=postgres \
--network rest-api-test \
--name rest-api-test appway/rest-api-test:latest
```

You can use these containers to test on the fly the application, and then you can remove everything afterwards. To
simply execute the **test** you can use the preconfigured script:

```shell
./test.sh
```

## Deploy and start the application (for prod env) [TO COMPLETE]

We assume for the production that you have:

- postgresql database running on some host

You can run the server as a docker container or compile it and executing on the server host.

### Run as a docker container [TODO]

...

### Run as a compiled binary [TODO]

...


