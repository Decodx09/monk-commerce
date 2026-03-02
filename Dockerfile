FROM golang:1.23-alpine

WORKDIR /app

# Pre-copy/cache go.mod for pre-downloading dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -o main .

EXPOSE 3000

CMD ["./main"]
