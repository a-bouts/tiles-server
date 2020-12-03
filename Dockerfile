FROM golang

WORKDIR /go/src/tiles-server

COPY . .

RUN go install -v ./...

ENTRYPOINT ["tiles-server"]
