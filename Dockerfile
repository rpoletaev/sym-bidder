FROM golang:latest

ADD . /go/src/github.com/rpoletaev/sym-bydder
WORKDIR /go/src/github.com/rpoletaev/sym-bydder
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bydder .

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/src/github.com/rpoletaev/sym-bydder .

CMD ["./bydder"]
EXPOSE 8080
