package gin

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"

	"gl.king.im/king-lib/framework/gin/server"

	kgin "github.com/go-kratos/gin"
	ke "github.com/go-kratos/kratos/v2/errors"
)

func AdaptToKratos(next middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {

		ginCtx, ok := kgin.FromGinContext(ctx)

		if !ok {
			return nil, ke.InternalServer("FROM_GIN_CONTEXT_NOT_OK", "")
		}

		//ctx = metadata.AppendToClientContext(ctx, framework.METADATA_KEY_CALL_SCENARIO, ginCtx.Request.Header.Get(framework.METADATA_KEY_CALL_SCENARIO))
		//ctx = metadata.AppendToClientContext(ctx, framework.METADATA_KEY_AUTH_TOKEN, ginCtx.Request.Header.Get(framework.METADATA_KEY_AUTH_TOKEN))

		//mdFromGin := metadata.New(map[string]string{
		//	framework.METADATA_KEY_CALL_SCENARIO: ginCtx.Request.Header.Get(framework.METADATA_KEY_CALL_SCENARIO),
		//	framework.METADATA_KEY_AUTH_TOKEN:    ginCtx.Request.Header.Get(framework.METADATA_KEY_AUTH_TOKEN),
		//})
		//
		//ctx = metadata.NewServerContext(ctx, mdFromGin)

		ctx = server.NewMetaServerCtxFromGin(ctx, ginCtx)

		return next(ctx, req)

	}
}
