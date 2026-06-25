# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependensi (di-cache oleh Docker jika go.mod tidak berubah)
COPY go.mod go.sum ./
RUN go mod download

# Install templ CLI
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy semua source code
COPY . .

# Generate file templ
RUN templ generate ./templates/...

# Build aplikasi
# - CGO_ENABLED=0 karena modernc.org/sqlite adalah pure Go
# - ldflags="-s -w" untuk memperkecil ukuran binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/activity-monitor ./main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install tzdata agar aplikasi bisa mengetahui timezone dengan benar (penting untuk reminder/deadline)
RUN apk --no-cache add tzdata

# Copy binary dari tahap builder
COPY --from=builder /app/bin/activity-monitor .

# Expose port aplikasi
EXPOSE 8080

# Command untuk menjalankan aplikasi
CMD ["./activity-monitor"]
