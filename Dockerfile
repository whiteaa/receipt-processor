FROM golang:1.22-alpine AS builder

WORKDIR /usr/src/app
COPY . .

RUN CGO_ENABLED=0 go build -o receipt-processor

FROM scratch

WORKDIR /app
EXPOSE 8080

COPY --from=builder /usr/src/app/receipt-processor .

CMD ["/app/receipt-processor"]