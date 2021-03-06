FROM golang:1.16.5-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/awsomeproject

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v ./...

# Build the Go app
RUN go build -o ./out/awsomeproject .


# Start fresh from a smaller image
FROM alpine:3.9 
ARG PORT=8080

RUN apk add ca-certificates

COPY --from=build_base /tmp/awsomeproject/out/awsomeproject /app/awsomeproject

# This container exposes port 8080 to the outside world
RUN echo $PORT
EXPOSE $PORT
ENV PORT=$PORT

# Run the binary program produced by `go install`
CMD ["/app/awsomeproject"]