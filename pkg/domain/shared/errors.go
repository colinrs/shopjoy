// Package shared 共享内核 - 错误定义
package shared

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/code"
)

// ErrInvalidTenantID 无效的租户ID
var ErrInvalidTenantID = &code.Err{
	HTTPCode: http.StatusBadRequest,
	Code:     200007,
	Msg:      "invalid tenant id",
}