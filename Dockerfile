FROM golang:alpine AS build

RUN apk update && apk --no-cache add ca-certificates dep git && update-ca-certificates

RUN mkdir /go/src/app
ADD . /go/src/app/
WORKDIR /go/src/app

RUN go build -o app main.go

FROM alpine AS runtime

COPY --from=build /go/src/app/* /
ENTRYPOINT ["/app"]
