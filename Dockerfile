# Stage 1: Build Docker image
FROM golang:1.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o .

# Final stage: Create minimal image with the built executable
FROM alpine:latest

COPY --from=builder /app/realchat /

EXPOSE 8080

CMD [ "/realchat" ]
