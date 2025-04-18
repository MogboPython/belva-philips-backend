FROM golang:1.24 AS builder
WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o main cmd/app/main.go

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=builder /app/main main
COPY /cmd/app/docs/* ./cmd/app/docs/
COPY /.env ./
COPY /internal/database/migrations ./internal/database/migrations


EXPOSE 8080

CMD ["/app/main"]