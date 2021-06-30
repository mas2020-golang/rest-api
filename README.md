# rest-api

Example project to build a rest API server with Go. The scope of this application is to use the smallest numbers of
extra packages to realize a full REST api server.

## Resources

This is the list of the reference to the resource that we are going to use in the application:

- **Gorilla Mux**, official documentation can be found [here](https://www.gorillatoolkit.org/). The Gorilla world on github
  is [here](https://github.com/gorilla).
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

## Start the application

To start the application, first set the correct environment variable:

```shell
export APP_DB_HOST=localhost \
export APP_DB_USERNAME=postgres \
export APP_DB_PASSWORD=password \
export APP_DB_NAME=postgres
```

then start a docker container as:

```shell
docker run  -d -p 5432:5432 -e POSTGRES_PASSWORD=password --name postgres_test postgres
```

finally execute the application typing:

```shell
go run *.go
```

## Test the application

To test, first add the environment variables, then execute:
```shell
go test -v github.com/mas2020-golang/rest-api/tests
```

## Curl examples for the 'products' handler

- **LOGIN** to the application

```shell
curl -v -s -X POST http://localhost:9090/login \
-H "Content-Type: multipart/form-data" \
-F 'username=andrea' -F 'password=test' | jq | \
sed 's+\([a-zA-Z0-9]*\.[a-zA-Z0-9]*\.[a-zA-Z0-9]*\).*+\1+'
```

- **GET** all the products

```shell
export TOKEN=
curl -v -s  http://localhost:9090/products \
-H "Authorization: Bearer ${TOKEN}" | jq
```

- **GET** the single product

```shell
curl -v -s  http://localhost:9090/products/1 \
-H "Authorization: Bearer {token}" | jq
```

- **CREATE** a new product

```shell
curl -s -X POST http://localhost:9090/products \
-H "Authorization: Bearer {token}" \
--data-binary @- << EOF | jq
{
    "name": "Espresso 2",
    "description": "Short and strong coffee",
    "price": 2.50,
    "sku": "dfadds-das-fdsa"
}
EOF
```

- **UPDATE** an existing product

```shell
curl -s -i -X PUT http://localhost:9090/products/1 \
-H "Authorization: Bearer {token}" \
--data-binary @- << EOF
{
    "name": "Espresso 900",
    "description": "More than a coffee",
    "price": 2.99,
    "sku": "df-d-fdsa"
}
EOF
```


