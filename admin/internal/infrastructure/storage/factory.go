package storage

import (
	"fmt"

	"github.com/colinrs/shopjoy/admin/internal/domain/media"
	snowflake "github.com/colinrs/shopjoy/pkg/snowflake"
)

// StorageType 存储类型
type StorageType string

const (
	StorageTypeLocal      StorageType = "local"
	StorageTypeCloudinary StorageType = "cloudinary"
	StorageTypeOSS        StorageType = "oss"
	StorageTypeS3         StorageType = "s3"
)

// Config 存储配置
type Config struct {
	Type       StorageType
	Local      LocalConfig
	Cloudinary CloudinaryConfig
	OSS        OSSConfig
	S3         S3Config
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	BasePath string
}

// CloudinaryConfig Cloudinary 配置
type CloudinaryConfig struct {
	CloudName    string
	APIKey       string
	APISecret    string
	UploadPreset string
	Environment  string
	Secure       bool
}

// OSSConfig OSS 配置
type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	Bucket          string
}

// S3Config S3 配置
type S3Config struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
}

// NewStorage factory.
func NewStorage(cfg Config, repo media.Repository, idGen snowflake.Snowflake) (Storage, error) {
	switch cfg.Type {
	case StorageTypeLocal:
		return newLocal(cfg, repo, idGen), nil
	case StorageTypeCloudinary:
		return newCloudinary(cfg, repo, idGen)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}

// MustNewStorage 创建存储实例，失败时panic
func MustNewStorage(cfg Config, repo media.Repository, idGen snowflake.Snowflake) Storage {
	s, err := NewStorage(cfg, repo, idGen)
	if err != nil {
		panic(err)
	}
	return s
}
