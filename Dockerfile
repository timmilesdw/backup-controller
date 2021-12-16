FROM golang:1.17.5-alpine3.15 as build

ARG name=backup-controller

WORKDIR /$name

COPY . .

RUN go mod download && go build -o ./bin/$name ./cmd/$name/main.go

FROM alpine:3.15

ARG name=backup-controller

USER 1001

COPY --from=build /$name/bin/$name .

ENTRYPOINT [ "/backup-controller" ]


