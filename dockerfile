FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./src/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

COPY /templates ./templates
COPY /static ./static

EXPOSE 8000
CMD ["./app"]