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
	MySQL infra.DBConfig
	JWT   JWTConfig
}
