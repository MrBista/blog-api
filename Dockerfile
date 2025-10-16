# Tahap 1: Build aplikasi
# Tetap menggunakan golang:1.25-alpine
FROM golang:1.25-alpine AS builder

# Set direktori kerja di dalam container
WORKDIR /app

# Salin go.mod dan go.sum
COPY go.mod .
COPY go.sum .

# Unduh dependensi (dependency caching)
RUN go mod download

# Salin kode sumber
COPY . .

# Build aplikasi Go
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /app/app ./cmd/main.go

# -----------------------------------------------------------------------------

# Tahap 2: Jalankan aplikasi (Menggunakan multi-stage build untuk image yang lebih kecil)
# GANTI: Menggunakan versi Alpine yang spesifik untuk keamanan dan konsistensi
FROM alpine:3.20 

# Instal sertifikat CA untuk permintaan HTTPS
RUN apk --no-cache add ca-certificates

# Set direktori kerja
WORKDIR /root/

# Salin binary yang sudah di-build dari tahap 'builder'
COPY --from=builder /app/app .
# !TODO ini di prod hapus
COPY config.yaml .

COPY wait-for-db.sh .
RUN chmod +x wait-for-db.sh

# Expose port yang digunakan Gofiber (biasanya 3000)
EXPOSE 3000

# Perintah default untuk menjalankan binary
CMD ["./wait-for-db.sh"]