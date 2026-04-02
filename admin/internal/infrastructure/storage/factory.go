package storage

import "fmt"

// StorageType 存储类型
type StorageType string

const (
	StorageTypeLocal StorageType = "local"
	StorageTypeOSS   StorageType = "oss"
	StorageTypeS3    StorageType = "s3"
)

// Config 存储配置
type Config struct {
	Type  StorageType
	Local LocalConfig
	OSS   OSSConfig
	S3    S3Config
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	BasePath string
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

// NewStorage 创建存储实例
func NewStorage(cfg Config) (Storage, error) {
	switch cfg.Type {
	case StorageTypeLocal:
		localCfg := cfg.Local
		if localCfg.BasePath == "" {
			localCfg.BasePath = "./uploads"
		}
		return &localStorage{
			basePath: localCfg.BasePath,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}

// MustNewStorage 创建存储实例，失败时panic
func MustNewStorage(cfg Config) Storage {
	s, err := NewStorage(cfg)
	if err != nil {
		panic(err)
	}
	return s
}
