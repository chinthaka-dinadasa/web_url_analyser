FROM golang:1.25.3-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download -x

COPY handlers models services ./

RUN go build -v -o main .

EXPOSE 8080

CMD ["./main"]