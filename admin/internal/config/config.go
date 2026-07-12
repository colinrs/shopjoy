package config

import (
	"github.com/colinrs/shopjoy/pkg/infra"
	"github.com/zeromicro/go-zero/rest"
)

type JWTConfig struct {
	Secret        string
	AccessExpiry  int64
	RefreshExpiry int64
}

type Config struct {
	rest.RestConf
	MySQL               infra.DBConfig
	JWT                 JWTConfig
	StripeWebhookSecret string
	Storage             StorageConf
}

// StorageConf 后端存储后端类型与配置。
type StorageConf struct {
	Type       string                `json:"type,default=local"`
	Local      LocalStorageConf      `json:"local,optional"`
	Cloudinary CloudinaryStorageConf `json:"cloudinary,optional"`
}

type LocalStorageConf struct {
	BasePath string `json:"base_path,default=./uploads"`
}

type CloudinaryStorageConf struct {
	CloudName    string `json:"cloud_name"`
	APIKey       string `json:"api_key"`
	APISecret    string `json:"api_secret"`
	UploadPreset string `json:"upload_preset,optional"`
	Environment  string `json:"environment,default=dev"`
	Secure       bool   `json:"secure,default=true"`
}
