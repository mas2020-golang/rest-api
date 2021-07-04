# STEP 1 build executable binary
FROM golang:1.16.5-buster AS builder

# Install git.
# Git is required for fetching the dependencies.
#RUN apk update && apk add --no-cache git

WORKDIR /usr/local/rest-api
COPY . .

# Using go get
RUN go get -d -v

# Build the binary (in case CGO_ENABLED=0 gives any problem remove it and ensure to use the same distro, e.g. alpine,
# for the builder the copy)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app *.go

# Step 2: copy the executable file
FROM alpine:latest
RUN mkdir -p /usr/local/rest-api/bin
COPY --from=builder /usr/local/rest-api/app /usr/local/rest-api/bin
WORKDIR /usr/local/rest-api/bin
EXPOSE 9090
CMD ["./app"]