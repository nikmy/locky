FROM golang:latest

COPY . /go/src/github.com/nikmy/locky

WORKDIR /go/src/github.com/nikmy/locky

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

STOPSIGNAL SIGINT

CMD ["/go/src/github.com/nikmy/locky/app"]