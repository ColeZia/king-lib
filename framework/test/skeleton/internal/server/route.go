package server

import (
	"github.com/gin-gonic/gin"

	"gl.king.im/king-lib/framework/service/boot/server"
	svcV1 "gl.king.im/king-lib/framework/test/skeleton/internal/service/service/v1"

	ginServer "gl.king.im/king-lib/framework/gin/server"
)

func NewEngine(skeleGinSvc *svcV1.SkeletonGinService) *gin.Engine {

	engine := gin.Default()

	//使用封装过的kratos中间件
	//engine.Use(kgin.Middlewares(fwGin.AdaptToKratos, fwGin.Base))
	engine.Use(server.GinBaseMiddlewares())

	RootG := engine.Group("/Root")
	{
		//模块A
		moduleA := RootG.Group("/moduleA")
		{
			moduleA.POST("/action1", func(gc *gin.Context) {
				ctx := ginServer.GinCtxToKratosCtx(gc)
				_ = ctx
				val, err := skeleGinSvc.GinTest(gc)
				ginServer.RspSerializes(gc, val, err)
			})
		}

		//模块B
		moduleB := RootG.Group("/moduleB")
		{
			moduleB.POST("/GinTest", func(gc *gin.Context) {
				ctx := ginServer.GinCtxToKratosCtx(gc)
				_ = ctx
				val, err := skeleGinSvc.GinTest(gc)
				ginServer.RspSerializes(gc, val, err)
			})

			moduleB.POST("/ErrorHandle", func(gc *gin.Context) {
				ctx := ginServer.GinCtxToKratosCtx(gc)
				_ = ctx
				_, err := skeleGinSvc.GinTest(gc)
				ginServer.ErrorHandle(gc, err, true)
			})

			moduleB.POST("/RspSerializes", func(gc *gin.Context) {
				ctx := ginServer.GinCtxToKratosCtx(gc)
				_ = ctx
				val, err := skeleGinSvc.GinTest(gc)
				ginServer.RspSerializes(gc, val, err)
			})
		}
	}

	return engine
}
