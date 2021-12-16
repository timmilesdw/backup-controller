FROM golang:1.17.5-alpine3.15 as build

ARG name=backup-controller

WORKDIR /$name

COPY go.mod .
COPY go.sum .

RUN go mod download 

COPY . .

RUN go build -o ./bin/$name ./cmd/$name/main.go

FROM alpine:3.15

ARG name=backup-controller

WORKDIR /$name

RUN apk update 
RUN apk add --no-cache postgresql-client

RUN chown -R 1001:1001 /$name
RUN chmod 755 /$name

USER 1001

COPY --from=build /$name/bin/$name .

ENTRYPOINT [ "/backup-controller/backup-controller" ]


