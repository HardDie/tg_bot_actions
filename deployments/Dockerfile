FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Build app
COPY . .
RUN go build -o bot ./cmd/main.go

# Parse criminal law records
RUN go run tools/parsing_criminal_law/main.go -url https://www.consultant.ru/document/cons_doc_LAW_10699/


# Build final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bot ./
COPY --from=builder /app/criminals.json ./

CMD ["/app/bot"]
