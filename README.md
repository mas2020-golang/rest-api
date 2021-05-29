# rest-api
Example project to build a rest API server with Go. The scope of this application is to use the smallest numbers of extra packages we can.

## Resources
This is the list of the reference to the resource that we are going to use in the application:
- Gorilla Mux, official documentation can be found [here](https://www.gorillatoolkit.org/).
  The Gorilla world on github is [here](https://github.com/gorilla).

## Structure of the application
The application has this folders:
- handlers: contains each http handler

## Curl examples for the 'product' handler 
- **GET** the product
```shell
curl -v -s  http://localhost:9090/products | jq
```
- **POST** the new product
```shell
curl -X POST -i  http://localhost:9090/products --data-binary @- << EOF   
{
    "id": 3,
    "name": "Espresso 2",
    "description": "Short and strong coffee",
    "price": 2.50,
    "sku": "dfadds"
}
EOF
```
- **PUT** the product
```shell
curl -i -X PUT http://localhost:9090/products/1 --data-binary @- << EOF   
{
    "id": 1,
    "name": "Espresso 900",
    "description": "More than a coffee",
    "price": 2.99,
    "sku": "h3218"
}
EOF
```


