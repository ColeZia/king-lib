package interceptors

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func BaseUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		//		// 把两个 ctx 合并成一个
		//		ctx, cancel := ic.Merge(ctx, s.ctx)
		//		defer cancel()
		//		// 从 ctx 中取出 metadata
		//		md, _ := grpcmd.FromIncomingContext(ctx)
		//		// 把一些信息绑定到 ctx 上
		//		ctx = transport.NewServerContext(ctx, &Transport{
		//			endpoint:  s.endpoint.String(),
		//			operation: info.FullMethod,
		//			header:    headerCarrier(md),
		//		})
		//		// ctx 超时设置
		//		if s.timeout > 0 {
		//			ctx, cancel = context.WithTimeout(ctx, s.timeout)
		//			defer cancel()
		//		}
		//		// 中间件处理
		//		h := func(ctx context.Context, req interface{}) (interface{}, error) {
		//			return handler(ctx, req)
		//		}
		//		if len(s.middleware) > 0 {
		//			h = middleware.Chain(s.middleware...)(h)
		//		}
		//		return h(ctx, req)
		log.Println("framework BaseUnaryInterceptor...")
		return handler(ctx, req)
	}
}
