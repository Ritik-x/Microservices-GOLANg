FROM golang:1.27-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o notification-service .

EXPOSE 8082

CMD ["./notification-service"]