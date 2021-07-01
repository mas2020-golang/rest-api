# STEP 1 build executable binary
FROM golang:1.16.5-buster

# Install git.
# Git is required for fetching the dependencies.
RUN apt-get update && apt-get install git

WORKDIR /usr/local/rest-api
COPY . .

# Using go get
RUN go get -d -v

# Execute test
CMD ["go", "test", "github.com/mas2020-golang/rest-api/test/..."]