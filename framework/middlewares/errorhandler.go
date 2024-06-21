package middlewares

import (
	"context"
	"log"

	ke "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

func ErrorHandler(next middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {

		rsp, err := next(ctx, req)
		se := ke.FromError(err)
		if se != nil {
			errCode := int(se.Code)
			if errCode >= 500 && errCode <= 599 {
				log.Println("ErrorHandler::", err)
				//err1 := errors.WithStack(err)
				//fmt.Printf("%+v", err1)
			}
		}
		return rsp, err
	}
}
