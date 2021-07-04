# Note

## Todo list

This paragraph contains a list of features we need, and we have not implemented yet in the project.

- **[TODO]** complete the documentation writing the section related to the production deployment
- **[TODO]** add the swagger engine to document API
- **[✔︎]** add the **test section** where to test every method before implementing it
- **[✔︎]** add **PostgreSQL database** to the project and change the model to read and write using a database connection
- **[✔︎]** implement the JWT security: add a POST /login resource path; create and check the JWT token
- **[✔︎]** add **JWT authorization** as a middleware layer (take a look first at https://github.com/pinbar/go-mux-jwt,
  the original library can be found [here](https://github.com/dgrijalva/jwt-go))
  To create the token use a random pwd that is a random UUID v4, for the package to use take a look [here](https://pkg.go.dev/github.com/google/uuid#section-readme).
- **[✔]** update the test code to include the token
- **[✔]** change the code to read username and password from a json body instead of a form
- **[✔]** add **goutils** to the project to better format the output
- **[✔]** add a Dockerfile.test to create an image for making the go test using a Postgres image on the same docker
  network. Then create a script .sh that starts the 2 containers and execute the test. (see if it is worth to use the
  official golang image in order to have an image to use for testing: now I fetch Alpine but I am not sure)
- **[✔]** introduce the Dockerfile to create a light image for the API server
- **[✔]** introduce the docker compose to have a fully working composition with rest-api and postgres
- **[✔]** add a **YAML configuration file** to the project