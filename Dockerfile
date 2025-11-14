FROM golang:1.25.3-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download -x

COPY main.go ./
COPY handlers/ ./handlers/
COPY models/ ./models/
COPY services/ ./services/
COPY logger/ ./logger/
COPY docs/ ./docs/

RUN go build -v -o main .

EXPOSE 8080

CMD ["./main"]