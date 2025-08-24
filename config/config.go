package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	JWT      JWTConfig      `mapstructure:"jwt"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Server   ServerConfig   `mapstructure:"server"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Env      string         `mapstructure:"env"`
	NATS     NATSConfig     `mapstructure:"nats"`
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	AccessTokenTTL  string `mapstructure:"access_token_ttl"`
	RefreshTokenTTL string `mapstructure:"refresh_token_ttl"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
}

// PostgresConfig holds PostgreSQL-specific configuration
type PostgresConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type NATSConfig struct {
	Brokers []string `mapstructure:"brokers"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
	IdleTimeout  string `mapstructure:"idle_timeout"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	BcryptCost     int    `mapstructure:"bcrypt_cost"`
	SessionTimeout string `mapstructure:"session_timeout"`
}

var (
	// Global config instance
	AppConfig *Config
)

// LoadConfig loads configuration from files based on environment
func LoadConfig() (*Config, error) {
	// Get environment from ENV variable, default to "develop"
	env := os.Getenv("ENV")
	if env == "" {
		env = "develop"
	}

	viper.SetConfigName("config") // base config file name
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// Read base config
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read base config: %w", err)
	}

	// Read environment-specific config
	viper.SetConfigName(env) // environment-specific config file name
	if err := viper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read %s config: %w", env, err)
	}

	// Parse configuration into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	AppConfig = &config
	return &config, nil
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	if config.Database.Postgres.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if config.Redis.Host == "" {
		return fmt.Errorf("redis host is required")
	}

	return nil
}

// GetPostgresDSN returns PostgreSQL connection string
func (c *Config) GetPostgresDSN() string {
	password, err := base64.StdEncoding.DecodeString(c.Database.Postgres.Password)
	if err != nil {
		return ""
	}
	pg := c.Database.Postgres
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pg.Host, pg.Port, pg.User, string(password), pg.DBName, pg.SSLMode)
}

// GetRedisAddr returns Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// GetAccessTokenTTL returns parsed access token TTL
func (c *Config) GetAccessTokenTTL() (time.Duration, error) {
	return time.ParseDuration(c.JWT.AccessTokenTTL)
}

// GetRefreshTokenTTL returns parsed refresh token TTL
func (c *Config) GetRefreshTokenTTL() (time.Duration, error) {
	return time.ParseDuration(c.JWT.RefreshTokenTTL)
}

// GetJWTSecret returns JWT secret as bytes
func (c *Config) GetJWTSecret() []byte {
	return []byte(c.JWT.Secret)
}

// Legacy functions for backward compatibility
func GetPostgresDSN() string {
	if AppConfig == nil {
		return "host=localhost user=root password=root dbname=authdb sslmode=disable"
	}
	return AppConfig.GetPostgresDSN()
}

func GetRedisAddr() string {
	if AppConfig == nil {
		return "localhost:6379"
	}
	return AppConfig.GetRedisAddr()
}

// Legacy variables for backward compatibility
var (
	JWTSecret        = []byte("Ym9va2luZy10aWNrZXQ=")
	AccessTokenTTL   = 15 * time.Minute
	RefreshTokenTTL  = 7 * 24 * time.Hour
	RedisTokenPrefix = "jwt_token:"
)
