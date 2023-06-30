FROM golang:1.19-alpine AS build

ARG APP_PORT

WORKDIR /usr/src/app/

COPY . /usr/src/app/

RUN go mod download

RUN go build -o /userservice

FROM alpine:latest

WORKDIR /

COPY --from=build /userservice /

EXPOSE ${APP_PORT}

CMD ["./userservice"]

