package service

import (
	"context"
	"gl.king.im/king-lib/framework/auth"
	"gl.king.im/king-lib/framework/internal/data"
	fwMd "gl.king.im/king-lib/framework/transport/metadata"
	authpb "gl.king.im/king-lib/protobuf/api/service/service/v1"
)

// 认证接口-服务间JWT实现
type ServiceInnerJWTAuthen struct {
}

var _ auth.Authentication = (*ServiceInnerJWTAuthen)(nil)

func (*ServiceInnerJWTAuthen) GetToken() (token string) {
	token = "ServiceInnerJWTAuthenToken"
	return token
}

func (*ServiceInnerJWTAuthen) Validate(ctx context.Context, token string) (ok bool, user interface{}, err error) {

	client, err := data.GetServiceServiceClient()
	if err != nil {
		return false, user, err
	}
	callCtx := fwMd.BuildInnerMDCtx()

	//	dataIns, err := data.OnceInitData()
	//	if err != nil {
	//		return false, err
	//	}
	//	client := dataIns.SrvCli
	//	callCtx := BuildMetaDataCtx()

	//	connGRPC, callCtx, err := NewGrpcClientConn("BossService")
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer connGRPC.Close()
	//
	//	client := authpb.NewServiceServiceClient(connGRPC)

	//登录认证判断
	reply, err := client.Authenticate(callCtx, &authpb.AuthenticateRequest{Token: token})
	if err == nil {
		return reply.Ok, user, nil
	} else {
		return false, user, err
	}

}

// 授权接口-服务间JWT实现
type ServiceInnerJWTAuthor struct {
}

var _ auth.Authorization = (*ServiceInnerJWTAuthor)(nil)

func (*ServiceInnerJWTAuthor) Can(ctx context.Context, token, resource string) (bool, error) {

	client, err := data.GetServiceServiceClient()
	if err != nil {
		return false, err
	}
	callCtx := fwMd.BuildInnerMDCtx()

	//	dataIns, err := data.OnceInitData()
	//	if err != nil {
	//		return false, err
	//	}
	//	client := dataIns.SrvCli
	//	callCtx := BuildMetaDataCtx()

	//	connGRPC, callCtx, err := NewGrpcClientConn("BossService")
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer connGRPC.Close()
	//
	//	client := authpb.NewServiceServiceClient(connGRPC)

	//权限判断
	canReply, err := client.AuthorizationCheck(callCtx, &authpb.AuthorizationCheckRequest{Token: token, Resource: resource})

	if err == nil {
		return canReply.Ok, nil
	} else {
		return false, err
	}

}

type ServiceInnerJWTAuth struct {
}

var _ auth.Auth = (*ServiceInnerJWTAuth)(nil)

func (*ServiceInnerJWTAuth) AuthCheck(ctx context.Context, token, resource string, aeIgnore, aoIgnore bool) (aeOk, aoOk bool, user interface{}, err error) {
	client, err := data.GetServiceServiceClient()
	if err != nil {
		return
	}

	//callCtx := fwMd.BuildInnerMDCtx()
	reply, err := client.AuthCheck(ctx, &authpb.AuthCheckRequest{Token: token, Resource: resource, AuthenticationIgnore: aeIgnore, AuthorizationIgnore: aoIgnore})
	if err != nil {
		return
	}

	user = reply.Service
	aeOk = reply.AuthenticateOk
	aoOk = reply.AuthorizationCheckOk

	return
}

type OpenapiAuth struct {
}

var _ auth.Auth = (*OpenapiAuth)(nil)

func (*OpenapiAuth) AuthCheck(ctx context.Context, token, resource string, aeIgnore, aoIgnore bool) (aeOk, aoOk bool, user interface{}, err error) {
	client, err := data.GetServiceServiceClient()
	if err != nil {
		return
	}

	callCtx := fwMd.BuildInnerMDCtx()
	reply, err := client.OpenapiAuthCheck(callCtx, &authpb.AuthCheckRequest{Token: token, Resource: resource, AuthenticationIgnore: aeIgnore, AuthorizationIgnore: aoIgnore})
	if err != nil {
		return
	}

	user = reply.Service
	aeOk = reply.AuthenticateOk
	aoOk = reply.AuthorizationCheckOk

	return
}
