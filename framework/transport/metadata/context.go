package metadata

import (
	"context"

	kmd "github.com/go-kratos/kratos/v2/metadata"
	"gl.king.im/king-lib/framework"
	"gl.king.im/king-lib/framework/auth/token"
	"gl.king.im/king-lib/framework/config"
	"gl.king.im/king-lib/framework/service/appinfo"
	grpcmd "google.golang.org/grpc/metadata"
)

func BuildInnerMDCtx() context.Context {
	callCtx := context.Background()
	callCtx = kmd.AppendToClientContext(callCtx, framework.METADATA_KEY_CALL_SCENARIO, framework.MDV_SERVICE_CALL_SCENARIO_INNER)

	//authToken := AppInfoIns.Authentication.GetToken()
	authToken := generateTokenWithConf()
	callCtx = kmd.AppendToClientContext(callCtx, framework.METADATA_KEY_AUTH_TOKEN, authToken)

	return callCtx
}

func BuildGatewayMDCtx(authToken string) context.Context {
	callCtx := context.Background()
	callCtx = kmd.AppendToClientContext(callCtx, framework.METADATA_KEY_CALL_SCENARIO, framework.MDV_SERVICE_CALL_SCENARIO_GETEWAY)

	//authToken := AppInfoIns.Authentication.GetToken()
	callCtx = kmd.AppendToClientContext(callCtx, framework.METADATA_KEY_AUTH_TOKEN, authToken)

	return callCtx
}

func BuildInnerMDCtxFromServerCtx(ctx context.Context) (cliCtx context.Context) {

	//authToken := AppInfoIns.Authentication.GetToken()
	authToken := generateTokenWithConf()

	if true {

		// x-md-global-
		//覆盖server里面的X-Md-Global-Call-Scenario和X-Md-Global-Auth-Token，解决下面提到的全局追踪id的传递又必须以server context为基础的问题
		//注意这里必须克隆metadata，否则直接覆盖的话，使用的是同一个map，会影响后续流程中读取server context的metadata，因为后续流程中认为server context的metadata是覆写之前的内容
		smd, _ := kmd.FromServerContext(ctx)
		newSmd := smd.Clone()
		//fmt.Println("BuildInnerMDCtxFromServerCtx smd, sok", sok, smd)

		newSmd.Set(framework.METADATA_KEY_CALL_SCENARIO, framework.MDV_SERVICE_CALL_SCENARIO_INNER)
		newSmd.Set(framework.METADATA_KEY_AUTH_TOKEN, authToken)

		//ctx = kmd.NewClientContext(ctx, nil)
		//由于前端的gateway header用的global的metadata，会被设置到client里面去，同时全局追踪id的传递又必须以server context为基础，所以，这里暂时强行将server的metadata置空，然后再设置服务间的auto token，此处可能会产生隐患，可能会把其他global md置空
		ctx = kmd.NewServerContext(ctx, newSmd)

		//旧版inner
		cliCtx = kmd.AppendToClientContext(ctx, framework.METADATA_KEY_CALL_SCENARIO, framework.MDV_SERVICE_CALL_SCENARIO_INNER)
		cliCtx = kmd.AppendToClientContext(cliCtx, framework.METADATA_KEY_AUTH_TOKEN, authToken)

		//新版
		//cliCtx = kmd.AppendToClientContext(cliCtx, framework.METADATA_KEY_LOCAL_SVC_TOKEN, authToken)

		//smd2, sok2 := metadata.FromServerContext(ctx)
		//fmt.Println("BuildInnerMDCtxFromServerCtx smd2, sok2", sok2, smd2)
	} else {
		keyvals := []string{
			framework.METADATA_KEY_CALL_SCENARIO,
			framework.MDV_SERVICE_CALL_SCENARIO_INNER,
			framework.METADATA_KEY_AUTH_TOKEN,
			authToken,
		}
		cliCtx = grpcmd.AppendToOutgoingContext(ctx, keyvals...)
	}

	return
}

func generateTokenWithConf() string {
	serviceConf := config.GetServiceConf()
	if serviceConf.Service.Secret == "" {
		panic("服务Secret未配置！")
	}

	token, err := token.GenerateToken(serviceConf.Service.Secret, appinfo.AppInfoIns.Name)

	if err != nil {
		panic("token生成失败！" + err.Error())
	}

	return token
}
