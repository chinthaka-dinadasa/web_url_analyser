FROM golang:1.25.3-alpine

WORKDIR /app

COPY vendor ./vendor

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN timeout 300 go build -mod=vendor -v -o main .

EXPOSE 8080

CMD ["./main"]