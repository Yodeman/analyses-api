# build stage
FROM golang:1.21-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o analyses-api main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

# run stage
FROM alpine:3.19
WORKDIR /app
RUN apk add make
COPY --from=builder /app/analyses-api .
COPY --from=builder /app/migrate /usr/bin/migrate
COPY start.sh .
COPY app.env .
COPY Makefile .
COPY dbase/migration ./dbase/migration

EXPOSE 8000
CMD [ "/app/analyses-api" ]
