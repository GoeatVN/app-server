# Sử dụng image Go chính thức
FROM golang:1.20-alpine AS builder

# Thiết lập thư mục làm việc
WORKDIR /app

# Copy các file go.mod và go.sum trước
COPY go.mod go.sum ./

# Tải các dependency
RUN go mod download

# Copy toàn bộ mã nguồn vào image
COPY . .

# Biên dịch mã nguồn
RUN go build -o main ./cmd/api/main.go

# Tạo một image nhỏ gọn để chạy ứng dụng
FROM alpine:3.18

# Thêm chứng chỉ SSL nếu cần kết nối với dịch vụ HTTPS
RUN apk --no-cache add ca-certificates

# Thiết lập thư mục làm việc
WORKDIR /root/

# Copy binary từ bước trước
COPY --from=builder /app/main .

# Expose cổng cho HTTP server
EXPOSE 8080

# Chạy ứng dụng
CMD ["./main"]
