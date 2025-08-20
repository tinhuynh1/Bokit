# Configuration System

Hệ thống cấu hình sử dụng Viper để load configuration từ file YAML dựa trên môi trường.

## Cấu trúc file

```
config/
├── config.yaml      # Base configuration
├── develop.yaml     # Development environment
├── uat.yaml         # UAT environment  
├── prod.yaml        # Production environment
├── config.go        # Configuration structs and loader
└── config_test.go   # Unit tests
```

## Cách sử dụng

### 1. Set Environment Variable

```bash
# Development (default)
export ENV=develop

# UAT
export ENV=uat

# Production
export ENV=prod
```

### 2. Load Configuration trong code

```go
import "booking-svc/config"

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Sử dụng configuration
    jwtSecret := cfg.GetJWTSecret()
    accessTTL, _ := cfg.GetAccessTokenTTL()
    dsn := cfg.GetPostgresDSN()
}
```

### 3. Cấu trúc Configuration

```go
type Config struct {
    JWT      JWTConfig      `mapstructure:"jwt"`
    Database DatabaseConfig  `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Server   ServerConfig   `mapstructure:"server"`
    Logging  LoggingConfig  `mapstructure:"logging"`
    Auth     AuthConfig     `mapstructure:"auth"`
}
```

## Các môi trường

### Development (develop.yaml)
- JWT secret: `dev-secret-key-123`
- Access token TTL: 30 phút
- Refresh token TTL: 24 giờ
- Database: localhost với SSL disable
- Redis: localhost không password
- Logging: debug level, text format

### UAT (uat.yaml)
- JWT secret: `uat-secret-key-456`
- Access token TTL: 15 phút
- Refresh token TTL: 7 ngày
- Database: UAT server với SSL require
- Redis: UAT server với password
- Logging: info level, json format

### Production (prod.yaml)
- JWT secret: `prod-secret-key-789`
- Access token TTL: 10 phút
- Refresh token TTL: 7 ngày
- Database: Production server với SSL require
- Redis: Production server với password
- Logging: warn level, json format

## Validation

Hệ thống tự động validate các giá trị bắt buộc:
- JWT secret không được rỗng
- Database host không được rỗng
- Redis host không được rỗng

## Legacy Support

Các hàm legacy vẫn được hỗ trợ để backward compatibility:

```go
// Legacy functions
dsn := config.GetPostgresDSN()
redisAddr := config.GetRedisAddr()
```

## Testing

Chạy test để kiểm tra configuration:

```bash
go test ./config -v
```

## Environment Variables

Có thể override các giá trị bằng environment variables:

```bash
export JWT_SECRET="custom-secret"
export DB_HOST="custom-host"
export REDIS_HOST="custom-redis"
```

## Security Notes

- Không commit file config chứa sensitive data vào git
- Sử dụng environment variables cho production secrets
- Rotate JWT secrets định kỳ
- Sử dụng strong passwords cho database và Redis 