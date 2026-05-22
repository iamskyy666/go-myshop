package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Top level config.
	Server ServerConfig
	Database DatabaseConfig
	JWT JWTConfig
	AWS AWSConfig
	Upload UploadConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret    string
	Expiry time.Duration
	RefreshTokenExpiry time.Duration
}

type AWSConfig struct {
	Region string
	AccessKeyID string
	SecretAccessKey string
	S3Bucket string
	S3Endpoint string
}

type UploadConfig struct {
	Path string
	MaxFileSize int64
}

func Load()(*Config, error){
	_ = godotenv.Load()

	jwtExpiry,_:=time.ParseDuration(getEnv("JWT_EXPIRY","24h"))
	refreshTokenExpiry,_:=time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY","72h"))
	maxUploadSize,_:=strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE","10485760"),10,64)

	return &Config{
		Server: ServerConfig{
	  	Port: getEnv("PORT","8080"),
		GinMode: getEnv("GIN_MODE","debug"),
		},
		Database: DatabaseConfig{
			Host: getEnv("DB_HOST","localhost"),
			Port: getEnv("DB_PORT","5432"),
			User: getEnv("DB_USER","postgres"),
			Password: getEnv("DB_PASSWORD","password"),
			Name: getEnv("DB_NAME","go-myshop"),
			SSLMode: getEnv("DB_SSL_MODE","disable"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET_KEY","my_secret_key@12345"),
			Expiry: jwtExpiry,
			RefreshTokenExpiry: refreshTokenExpiry,
		},
		AWS: AWSConfig{
			Region: getEnv("AWS_REGION","us-east-1"),
			AccessKeyID: getEnv("AWS_ACCESS_KEY_ID","test"),
		    SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY","test"),
			S3Bucket: getEnv("AWS_S3_BUCKET","go-myshop-uploads"),
			S3Endpoint: getEnv("AWS_S3_ENDPOINT","http://localhost:4566"),
		},
		Upload: UploadConfig{
			Path: getEnv("UPLOAD_PATH","./uploads"),
			MaxFileSize: maxUploadSize,
		},
	},nil
}

func getEnv(key,defaultVal string)string{
if val:=os.Getenv(key);val!=""{
	return val
}
	return defaultVal
}