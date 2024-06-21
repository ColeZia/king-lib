package auth

import (
	"context"
	"io"
	"log"
	"net/http"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"gl.king.im/king-lib/framework/auth"
	"gl.king.im/king-lib/framework/internal/data"
	"gl.king.im/king-lib/framework/service"
	"gl.king.im/king-lib/framework/transport/http/client"
	hc "gl.king.im/king-lib/framework/transport/http/client"
	fwMd "gl.king.im/king-lib/framework/transport/metadata"
	"google.golang.org/protobuf/encoding/protojson"

	ke "github.com/go-kratos/kratos/v2/errors"
	"gl.king.im/king-lib/protobuf/api/common"
	svcPb "gl.king.im/king-lib/protobuf/api/service/service/v1"
	userPb "gl.king.im/king-lib/protobuf/api/user/service/v1"
	"google.golang.org/protobuf/proto"
)

func AuthTokenDefault(ctx context.Context, authToken, op string, authIgnore, authenticationIgnore, authorizationIgnore bool, caller string) (err error, user interface{}) {
	if authIgnore {
		return nil, user
	}
	//判断服务是否有注册对应的自定义鉴权方案，如果没有则报错
	val, ok := service.AppInfoIns.ScenarioAuthRegistry.Load(caller)
	if ok {

		ar := val.(auth.AuthRegister)
		//认证
		if !authenticationIgnore {
			ok, user, err := ar.Ae.Validate(ctx, authToken)
			if err != nil {
				return err, user
			}

			if !ok {
				kerr := ke.New(401, "BASE_MW_"+caller+"_AUTHENTICATION_ERROR", "当前未认证")
				return kerr, user
			}
		}
		//权限
		if !authorizationIgnore {
			ok, err := ar.Ao.Can(ctx, authToken, op)
			if err != nil {
				return err, user
			}

			if !ok {
				return ke.New(403, "BASE_MW_"+caller+"_AUTHORIZATION_ERROR", "此操作未授权"), user
			}
		}
	} else {
		kerr := ke.New(400, "BASE_MW_CALLER_INFO_ERROR", "不支持的认证方式！")
		return kerr, user
	}
	return nil, user
}

func AuthTokenInner(ctx context.Context, authToken, op string, authIgnore, authenticationIgnore, authorizationIgnore bool) (err error, user interface{}) {
	if authIgnore {
		return
	}

	//如果是服务间的调用-则需要做认证和鉴权
	//在validate和can接口被调用的时候不用再发起validate和can认证授权校验了，否则就是进入无限循环了
	if op == "/api.auth.v1.ServiceAuthentication/Validate" || op == "/api.auth.v1.ServiceAuthorization/Can" {
		return
	}

	//如果是gateway的认证授权接口也不用做服务间的认证鉴权了，因为这个是比较基础的接口，每一个服务都赋权的话会比较麻烦，或者待后续创建一个权限角色之后统一赋权
	if op == "/api.auth.v1.Auth/IsLogin" || op == "/api.auth.v1.Auth/Can" {
		return
	}

	if authToken == "" {
		kerr := ke.New(401, "BASE_MW_SERVICE_AUTH_TOKEN_ERROR", "AUTH_TOKEN缺失！")
		return kerr, user
	}

	//暂不开放自定义认证鉴权
	//if service.AppInfoIns.Authentication == nil {
	if true {
		//认证+授权校验
		au := service.ServiceInnerJWTAuth{}
		var aeOk, aoOk bool

		aeOk, aoOk, user, err = au.AuthCheck(ctx, authToken, op, authenticationIgnore, authorizationIgnore)

		if err != nil {
			return
		}

		if !aeOk {
			err = ke.New(401, "BASE_MW_SERVICE_AUTHENTICATION_ERROR", "认证失败！")
			return
		}

		if !aoOk {
			err = ke.New(403, "BASE_MW_SERVICE_AUTHORIZATION_ERROR", "未授权")
			return
		}
		//		//如果微服务自己未注册认证和鉴权功能，则调用统一的认证鉴权
		//		//认证
		//		if !authenticationIgnore {
		//			ae := service.ServiceInnerJWTAuthen{}
		//			ok, user, err := ae.Validate(ctx, authToken)
		//			if err != nil {
		//				return err, user
		//			}
		//
		//			if !ok {
		//				fmt.Println("BASE_MW_SERVICE_AUTHENTICATION_ERROR op::", op)
		//
		//				kerr := ke.New(401, "BASE_MW_SERVICE_AUTHENTICATION_ERROR", "认证失败！")
		//
		//				return kerr, user
		//			}
		//		}
		//
		//		//授权
		//		if !authorizationIgnore {
		//			ao := service.ServiceInnerJWTAuthor{}
		//			ok, err := ao.Can(ctx, authToken, op)
		//			if err != nil {
		//				return err, user
		//			}
		//
		//			if !ok {
		//				fmt.Println("BASE_MW_SERVICE_AUTHORIZATION_ERROR op::", op)
		//
		//				kerr := ke.New(403, "BASE_MW_SERVICE_AUTHORIZATION_ERROR", "未授权")
		//				return kerr, user
		//			}
		//		}

	} else {
		//否则调用注册的实现
		if !authenticationIgnore {
			ok, user, err := service.AppInfoIns.Authentication.Validate(ctx, authToken)
			if err != nil {
				return err, user
			}

			if !ok {
				kerr := ke.New(401, "BASE_MW_SERVICE_AUTHENTICATION_ERROR", "当前未认证")
				return kerr, user
			}
		}

		if !authorizationIgnore {
			ok, err := service.AppInfoIns.Authorization.Can(ctx, authToken, op)
			if err != nil {
				return err, user
			}

			if !ok {
				return ke.New(403, "BASE_MW_SERVICE_AUTHORIZATION_ERROR", "此操作未授权"), user
			}
		}
	}

	return nil, user
}

func AuthTokenGatewayScenario(ctx context.Context, authToken, op string, authIgnore, authenticationIgnore, authorizationIgnore bool) (err error, user interface{}) {
	if authIgnore {
		return
	}

	logFlag := "AuthTokenGatewayScenario::"
	//在认证和授权校验接口被调用的时候不用再发起validate和can认证授权校验了，否则就是进入无限循环了
	if op == "/api.auth.v1.Auth/Login" || op == "/api.auth.v1.Auth/IsLogin" {
		return
	}

	if authToken == "" {
		err = ke.New(401, "BASE_MW_GATEWAY_AUTH_TOKEN_ERROR", "AUTH_TOKEN缺失！")
		return
	}

	client, err := data.GetUserServiceClient()
	if err != nil {
		return err, user
	}
	//callCtx := fwMd.BuildInnerMDCtx()

	//认证+授权校验
	var reply *userPb.AuthCheckReply
	reply, err = client.AuthCheck(ctx, &userPb.AuthCheckRequest{Token: authToken,
		Resource:             op,
		AuthenticationIgnore: authenticationIgnore,
		AuthorizationIgnore:  authorizationIgnore})

	if err != nil {
		log.Println(logFlag, "client.AuthCheck err::", err)
		return
	}

	if !reply.AuthenticateOk {
		log.Println(logFlag, "BASE_MW_GATEWAY_SCEN_AE_ERROR op::", op)
		err = ke.New(401, "BASE_MW_GATEWAY_SCEN_AE_ERROR", "当前未登录")
		return
	}

	if !reply.AuthorizationCheckOk {
		log.Println(logFlag, "BASE_MW_GATEWAY_SCEN_AO_ERROR op::", op)
		err = ke.New(403, "BASE_MW_GATEWAY_SCEN_AO_ERROR", "此操作未授权")
		return
	}

	user = reply.User

	//登录认证判断
	//	if !authenticationIgnore {
	//		reply, err := client.Authenticate(callCtx, &userPb.AuthenticateRequest{Token: authToken})
	//		if err != nil {
	//			return err, user
	//		}
	//
	//		if !reply.Ok {
	//			log.Println(logFlag, "BASE_MW_GATEWAY_SCEN_GTW_AE_ERROR op::", op)
	//			kerr := ke.New(401, "BASE_MW_GATEWAY_SCEN_GTW_AE_ERROR", "当前未登录")
	//			return kerr, user
	//		}
	//	}
	//
	//	//权限判断
	//	if !authorizationIgnore {
	//		canReply, err := client.AuthorizationCheck(callCtx, &userPb.AuthorizationCheckRequest{Token: authToken, Resource: op})
	//		if err != nil {
	//			return err, user
	//		}
	//
	//		if !canReply.Ok {
	//			log.Println(logFlag, "BASE_MW_GATEWAY_SCEN_GTW_AO_ERROR op::", op)
	//			kerr := ke.New(403, "BASE_MW_GATEWAY_SCEN_GTW_AO_ERROR", "此操作未授权")
	//			return kerr, user
	//
	//		}
	//	}
	return
}

func OpenapiAuth(ctx context.Context, authToken, op string, authIgnore, authenticationIgnore, authorizationIgnore bool) (err error, user interface{}) {
	if authIgnore {
		return
	}

	if authToken == "" {
		kerr := ke.New(401, "BASE_MW_OPENAPI_AUTH_TOKEN_ERROR", "AUTH_TOKEN缺失！")
		return kerr, user
	}

	//暂不开放自定义认证鉴权
	//if service.AppInfoIns.Authentication == nil {
	if true {
		//认证+授权校验
		au := service.OpenapiAuth{}
		var aeOk, aoOk bool
		aeOk, aoOk, user, err = au.AuthCheck(ctx, authToken, op, authenticationIgnore, authorizationIgnore)

		if err != nil {
			return
		}

		if !aeOk {
			err = ke.New(401, "BASE_MW_SERVICE_AUTHENTICATION_ERROR", "认证失败！")
			return
		}

		if !aoOk {
			err = ke.New(403, "BASE_MW_SERVICE_AUTHORIZATION_ERROR", "未授权")
			return
		}

	}

	return nil, user
}

func InnerAuthCheckInExternalNetwork(ctx context.Context, authToken, op string, authIgnore, authenticationIgnore, authorizationIgnore bool) (err error, user interface{}) {
	conn, callCtx, err := client.NewHttpClientConnInExtNet(common.SERVICE_NAME_BossService.String(), khttp.WithResponseDecoder(hc.JsonpbDecoder))
	if err != nil {
		return
	}
	_ = callCtx

	args := &svcPb.AuthCheckRequest{
		Token:                authToken,
		Resource:             op,
		AuthenticationIgnore: authenticationIgnore,
		AuthorizationIgnore:  authorizationIgnore,
	}

	reply := &svcPb.AuthCheckReply{}

	//测试
	//if false {
	//	gogoProto.RegisterEnum("common.BossService_Status", common.BossService_Status_name, common.BossService_Status_value)
	//	gogoProto.RegisterEnum("common.SERVICE_PLATFORM_TYPE", common.SERVICE_PLATFORM_TYPE_name, common.SERVICE_PLATFORM_TYPE_value)
	//}

	//	if false {
	//		um := jsonpb.Unmarshaler{}
	//		outStr := `{
	//			"authenticate_ok": true,
	//			"authorization_check_ok": true,
	//			"service": {
	//				"access_key": "",
	//				"name": "BossFrameworkTestpppp",
	//				"platform": "BOSS",
	//				"secret": "",
	//				"status": 1
	//			}
	//		}`
	//
	//		err = um.Unmarshal(strings.NewReader(string(outStr)), reply)
	//
	//		log.Println(err)
	//	}

	err = conn.InvokePBWithPathHandle(ctx, "POST", "/ServiceService/v1/AuthCheckForExternalNetwork", args, reply)

	if err != nil {
		log.Println("InnerAuthCheckInExternalNetwork InvokePBWithPathHandle() err::", err)
		return
	}

	user = reply.Service
	aeOk := reply.AuthenticateOk
	aoOk := reply.AuthorizationCheckOk

	if err != nil {
		return
	}

	if !aeOk {
		err = ke.New(401, "BASE_MW_SERVICE_AUTHENTICATION_ERROR", "认证失败！")
		return
	}

	if !aoOk {
		err = ke.New(403, "BASE_MW_SERVICE_AUTHORIZATION_ERROR", "未授权")
		return
	}

	return

}

var DecodeResponseFunc = func(ctx context.Context, res *http.Response, out interface{}) (err error) {
	//um := protojson.UnmarshalOptions{}
	//err = um.Unmarshal(res)

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	um := protojson.UnmarshalOptions{}
	err = um.Unmarshal(data, out.(proto.Message))

	return
}

func GatewayScenarioAuthCheckInExternalNetwork(ctx context.Context, authToken, op string, authIgnore, authenticationIgnore, authorizationIgnore bool) (err error, user *common.BossOperationUser) {
	if authIgnore {
		return
	}

	logFlag := "AuthTokenGatewayScenario::"

	if authToken == "" {
		err = ke.New(401, "BASE_MW_GATEWAY_AUTH_TOKEN_ERROR", "AUTH_TOKEN缺失！")
		return
	}

	//conn, callCtx, err := service.NewHttpClientConnInExtNet(common.SERVICE_NAME_BossUser.String(), khttp.WithResponseDecoder(DecodeResponseFunc))
	conn, callCtx, err := client.NewHttpClientConnInExtNet(common.SERVICE_NAME_BossUser.String())
	if err != nil {
		return
	}
	_ = callCtx

	args := &userPb.AuthCheckRequest{Token: authToken,
		Resource:             op,
		AuthenticationIgnore: authenticationIgnore,
		AuthorizationIgnore:  authorizationIgnore}

	callCtx1 := fwMd.BuildInnerMDCtx()

	//type commonWrapper struct {
	//	Code      int
	//	Data      interface{}
	//	Ts        string
	//	Message   string
	//	RequestId string
	//}

	type HttpBossOperationUser struct {
		Id             uint64               `protobuf:"varint,100,opt,name=id,proto3" json:"id,string,omitempty"`
		Username       string               `protobuf:"bytes,200,opt,name=username,proto3" json:"username,omitempty"`
		Nickname       string               `protobuf:"bytes,300,opt,name=nickname,proto3" json:"nickname,omitempty"`
		Name           string               `protobuf:"bytes,400,opt,name=name,proto3" json:"name,omitempty"` //姓名
		Mobile         string               `protobuf:"bytes,500,opt,name=mobile,proto3" json:"mobile,omitempty"`
		Gender         string               `protobuf:"varint,600,opt,name=gender,proto3,enum=common.BossOperationUser_Gender" json:"gender,omitempty"`
		Email          string               `protobuf:"bytes,700,opt,name=email,proto3" json:"email,omitempty"`   //邮箱
		Avatar         string               `protobuf:"bytes,800,opt,name=avatar,proto3" json:"avatar,omitempty"` //头像
		Employee       *common.BossEmployee `protobuf:"bytes,900,opt,name=employee,proto3" json:"employee,omitempty"`
		OceanAuthToken string               `protobuf:"bytes,1000,opt,name=ocean_auth_token,json=oceanAuthToken,proto3" json:"ocean_auth_token,omitempty"`
	}

	type HttpAuthCheckReply struct {
		AuthenticateOk       bool                   `protobuf:"varint,1,opt,name=authenticate_ok,json=authenticateOk,proto3" json:"authenticate_ok,omitempty"`
		AuthorizationCheckOk bool                   `protobuf:"varint,2,opt,name=authorization_check_ok,json=authorizationCheckOk,proto3" json:"authorization_check_ok,omitempty"`
		User                 *HttpBossOperationUser `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
	}

	//replywrapper := &commonWrapper{
	//	Data: &HttpAuthCheckReply{},
	//}

	reply := &HttpAuthCheckReply{}

	err = conn.InvokeWithPathHandle(callCtx1, "POST", "/UserService/v1/AuthCheckForExternalNetwork", args, &reply)
	//err = conn.Invoke(callCtx1, "POST", "/UserService/v1/AuthCheck", args, &replywrapper)

	//reply := replywrapper.Data.(*HttpAuthCheckReply)
	//认证+授权校验

	if err != nil {
		log.Println(logFlag, "client.AuthCheck err::", err)
		return
	}

	if !reply.AuthenticateOk {
		log.Println(logFlag, "BASE_MW_GATEWAY_SCEN_AE_ERROR op::", op)
		err = ke.New(401, "BASE_MW_GATEWAY_SCEN_AE_ERROR", "当前未登录")
		return
	}

	if !reply.AuthorizationCheckOk {
		log.Println(logFlag, "BASE_MW_GATEWAY_SCEN_AO_ERROR op::", op)
		err = ke.New(403, "BASE_MW_GATEWAY_SCEN_AO_ERROR", "此操作未授权")
		return
	}

	user = &common.BossOperationUser{
		Id:             reply.User.Id,
		Username:       reply.User.Username,
		Nickname:       reply.User.Nickname,
		Name:           reply.User.Name,
		Mobile:         reply.User.Mobile,
		Email:          reply.User.Email,
		Avatar:         reply.User.Avatar,
		OceanAuthToken: reply.User.OceanAuthToken,
	}

	if reply.User.Gender != "" {
		if val, ok := common.BossOperationUser_Gender_value[reply.User.Gender]; ok {
			user.Gender = common.BossOperationUser_Gender(val)
		}

	}
	//user = reply.User

	return
}
