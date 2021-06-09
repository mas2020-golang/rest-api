# rest-api

Example project to build a rest API server with Go. The scope of this application is to use the smallest numbers of
extra packages to realize a full REST api server.

## Resources

This is the list of the reference to the resource that we are going to use in the application:

- Gorilla Mux, official documentation can be found [here](https://www.gorillatoolkit.org/). The Gorilla world on github
  is [here](https://github.com/gorilla).
- Validation of JSON object is done using the package ***validator***. The repository can be
  found [here](https://github.com/go-playground/validator).
- To connect to PostgreSQL has been used the [pgx](https://pkg.go.dev/github.com/jackc/pgx) package

For the principles to follow when creating a REST API application take a look at
this [article](https://docs.microsoft.com/en-us/azure/architecture/best-practices/api-design).

## Remaining features to add

This paragraph contains a list of features we need, and we have not implemented yet in the project.

- add **PostgreSQL database** to the project and change the model to read and write using a database connection
- add **authentication** to the project with a POST /login method
- add **JWT authorization** as a middleware layer
- add **goutils** to the project to better format the output
- add a **YAML configuration file** to the project
- add the **test section** where to test every method before implementing it

## Dependencies

Before to start ensure to have the correct packages:

```shell
go get github.com/jackc/pgx/v4
go get github.com/jackc/pgx/v4/pgxpool
```

## Structure of the application

The application has these folders:

- handlers: contains each http handler

## Start the application

To start the application, first set the correct environment variable:

```shell
To connect for testing using Docker execute:
```shell
export APP_DB_HOST=localhost \
export APP_DB_USERNAME=postgres \
export APP_DB_PASSWORD=password \
export APP_DB_NAME=postgres
```

the start a docker container as:

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
go test -v
```

## Curl examples for the 'product' handler

- **GET** all the products

```shell
curl -v -s  http://localhost:9090/products | jq
```

- **GET** the single product

```shell
curl -v -s  http://localhost:9090/products/1 | jq
```

- **POST** the new product

```shell
curl -X POST -i  http://localhost:9090/products --data-binary @- << EOF   
{
    "name": "Espresso 2",
    "description": "Short and strong coffee",
    "price": 2.50,
    "sku": "dfadds-das-fdsa"
}
EOF
```

- **PUT** the product

```shell
curl -i -X PUT http://localhost:9090/products/1 --data-binary @- << EOF   
{
    "name": "Espresso 900",
    "description": "More than a coffee",
    "price": 2.99,
    "sku": "df-d-fdsa"
}
EOF
```


