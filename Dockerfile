FROM golang:alpine as builder

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY ./main.go ./
COPY ./database ./database/
COPY ./handler ./handler/
COPY ./socket ./socket/
COPY ./util ./util/

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build \
    -o /go/bin/main \
    -ldflags '-s -w'

EXPOSE 3000

FROM scratch as runner

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/main /app/main

ENTRYPOINT ["/app/main"]
