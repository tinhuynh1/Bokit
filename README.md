# 🎫 Event Booking Service - Hướng Dẫn Sử Dụng & Testing

## 📋 Mục Lục

- [Tổng Quan](#-tổng-quan)
- [Cài Đặt](#-cài-đặt)
- [Khởi Chạy](#-khởi-chạy)
- [API Documentation](#-api-documentation)
- [Testing](#-testing)
- [Troubleshooting](#-troubleshooting)

## 🎯 Tổng Quan

Event Booking Service là một hệ thống quản lý sự kiện và đặt vé được xây dựng bằng Go, cung cấp:

- **Quản lý sự kiện** (tạo, cập nhật, xóa)
- **Đặt vé** với kiểm tra số lượng còn lại
- **Xử lý thanh toán** qua NATS messaging
- **Thống kê** doanh thu và số vé đã bán
- **Authentication** JWT với role-based access
- **Scheduled jobs** tự động hủy booking hết hạn

## 🛠️ Cài Đặt

### Yêu Cầu Hệ Thống

- **Go 1.23+**
- **Docker & Docker Compose**
- **Goose** (database migration tool)

### Cài Đặt Go

**macOS:**
```bash
brew install go
go version
```

**Ubuntu/Debian:**
```bash
wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

**Windows:**
- Tải từ [https://go.dev/dl/](https://go.dev/dl/)
- Chạy installer và làm theo hướng dẫn

### Cài Đặt Docker

**macOS:**
```bash
brew install --cask docker
open /Applications/Docker.app
docker --version
```

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install apt-transport-https ca-certificates curl gnupg lsb-release
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin
sudo systemctl start docker
sudo usermod -aG docker $USER
docker --version
```

### Cài Đặt Goose

```bash
# Sử dụng Go
go install github.com/pressly/goose/v3/cmd/goose@latest

# Hoặc macOS
brew install goose

# Kiểm tra
goose --version
```

## 🚀 Khởi Chạy

### 1. Clone Repository

```bash
git clone <your-repo-url>
cd auth-service
```

### 2. Khởi Động Dependencies

```bash
# Khởi động tất cả services
docker-compose up -d

# Kiểm tra trạng thái
docker-compose ps
```

**Services được khởi động:**
- **PostgreSQL** (port 5432) - Database chính
- **Redis** (port 6379) - Cache và distributed lock
- **NATS** (port 4222) - Message broker
- **NATS UI** (port 8080) - Giao diện quản lý NATS

### 3. Chạy Database Migrations

```bash
# Chạy migrations
make migrate-up

# Hoặc thủ công
goose -dir migration postgres "host=localhost port=5432 user=booking_app password=X4pV7_qM9%tN1wK6@rG8jM2Z dbname=booking sslmode=disable" up
```

### 4. Khởi Chạy Service

```bash
# Sử dụng Make
make run

# Hoặc trực tiếp
go run main.go
```

Service sẽ chạy tại `http://localhost:8080`

## 📚 API Documentation

### Authentication

Tất cả API endpoints (trừ `/health`) yêu cầu JWT token trong header:

```bash
Authorization: Bearer <your-jwt-token>
```

### Sample Tokens (cho testing)

**Admin Token:**


