FROM golang:1.20 AS build
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o simple-proxy


FROM alpine
WORKDIR /app
COPY --from=build /go/src/app/simple-proxy .
