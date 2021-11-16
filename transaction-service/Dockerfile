FROM golang:alpine AS build-env
ENV GOPATH /go
WORKDIR /go/src
COPY . /go/src/transaction-service
RUN cd /go/src/transaction-service && go build .

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk*
WORKDIR /app
COPY --from=build-env /go/src/transaction-service/transaction-service /app

ENTRYPOINT [ "./transaction-service" ]