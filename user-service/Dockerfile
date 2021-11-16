FROM golang:alpine AS build-env
ENV GOPATH /go
WORKDIR /go/src
COPY . /go/src/user-service
RUN cd /go/src/user-service && go build .

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk*
WORKDIR /app
COPY --from=build-env /go/src/user-service/user-service /app

ENTRYPOINT [ "./user-service" ]