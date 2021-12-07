FROM golang:alpine

WORKDIR /go/src/app

COPY . /go/src/app

RUN go get -u ./...

CMD ["go","run","main.go"]