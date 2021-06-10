# Note

## Remaining features to add

This paragraph contains a list of features we need, and we have not implemented yet in the project.

- **[✔︎]** add **PostgreSQL database** to the project and change the model to read and write using a database connection
- **[✔︎]** implement the JWT security: add a POST /login resource path; create and check the JWT token
- **[TODO]** add **JWT authorization** as a middleware layer (take a look first at https://github.com/pinbar/go-mux-jwt,
  the original library can be found [here](https://github.com/dgrijalva/jwt-go))
  To create the token use a random pwd that is a random UUID v4, for the package to use take a look [here](https://pkg.go.dev/github.com/google/uuid#section-readme).
- **[TODO]** add **goutils** to the project to better format the output
- **[TODO]** add a **YAML configuration file** to the project
- **[TODO]** review the test part to ensure a dedicated container only for testing purposes (as done in the postgresql go lab)
- **[✔︎]** add the **test section** where to test every method before implementing it
