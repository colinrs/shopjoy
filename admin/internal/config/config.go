package config

import (
	"github.com/colinrs/shopjoy/pkg/infra"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	MySQL infra.DBConfig
}
