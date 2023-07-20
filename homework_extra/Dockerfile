FROM golang:1.19

ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o app cmd/server/main.go

CMD ["./app", "-adrepo"]