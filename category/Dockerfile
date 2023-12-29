FROM golang:1.21.5-alpine AS build
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /go/bin/app

FROM alpine:3.17

WORKDIR /

COPY --from=build /go/bin/app /app

EXPOSE 8080
RUN adduser -u 1001 -D -s /bin/sh -g ping 1001
USER 1001:1001

ENTRYPOINT ["/app"]