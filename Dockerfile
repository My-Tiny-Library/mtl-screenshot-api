FROM golang:latest as build

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

FROM chromedp/headless-shell:latest

RUN apt-get update
RUN apt-get install tini -y

WORKDIR /dist
COPY --from=build /build/main /dist/main

ENTRYPOINT ["tini", "--"]

EXPOSE 3000

CMD ["/dist/main"]