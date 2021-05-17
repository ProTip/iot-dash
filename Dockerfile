FROM golang:1.16

WORKDIR /go/src/app
COPY server .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["go", "run", "."]