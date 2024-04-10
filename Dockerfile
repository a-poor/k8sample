FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /app/k8sample-server .

FROM alpine
WORKDIR /app
COPY --from=builder /app/k8sample-server /k8sample-server
ENV APP_HOST="0.0.0.0"
ENV APP_PORT=8080
EXPOSE 8080
ENTRYPOINT [ "/app/k8sample-server" ]

