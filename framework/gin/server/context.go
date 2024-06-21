package server

import (
	"context"

	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/metadata"
	"gl.king.im/king-lib/framework"
)

func GinCtxToKratosCtx(gc *gin.Context) context.Context {
	// NewGinContext returns a new Context that carries gin.Context value.
	ctx := kgin.NewGinContext(gc.Request.Context(), gc)

	return NewMetaServerCtxFromGin(ctx, gc)
}

func NewMetaServerCtxFromGin(ctx context.Context, gc *gin.Context) context.Context {
	mdFromGin := metadata.New(map[string]string{
		framework.METADATA_KEY_CALL_SCENARIO: gc.Request.Header.Get(framework.METADATA_KEY_CALL_SCENARIO),
		framework.METADATA_KEY_AUTH_TOKEN:    gc.Request.Header.Get(framework.METADATA_KEY_AUTH_TOKEN),
	})

	ctx = metadata.NewServerContext(ctx, mdFromGin)

	return ctx
}
