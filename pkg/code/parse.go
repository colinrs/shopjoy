package code

import (
	"context"
	"errors"

	"github.com/spf13/cast"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func ParseErrToCodeMsg(_ context.Context, err error) *Err {
	if err == nil {
		return OKErr
	}
	rpcErr, ok := status.FromError(err)
	if ok {
		return NewErr(WithCode(cast.ToInt(rpcErr.Code())), WithMsg(rpcErr.String()))
	}
	var codeMsg *Err
	switch {
	case errors.As(err, &codeMsg):
	case errors.Is(err, gorm.ErrRecordNotFound):
		codeMsg = ErrNotFound.Copy().SetMsg(err.Error())
	default:
		codeMsg = UnknownErr.Copy().SetMsg(err.Error())
	}
	return codeMsg
}
