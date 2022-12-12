FROM golang:alpine as build-env

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf

WORKDIR /scrubber

COPY . .

RUN go build -tags musl -o scrubber .

CMD ./scrubber