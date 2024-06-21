package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"gl.king.im/king-lib/framework"
	"gl.king.im/king-lib/framework/auth"
	authUser "gl.king.im/king-lib/framework/auth/user"
	"gl.king.im/king-lib/framework/config"
	"gl.king.im/king-lib/framework/internal/data"
	"gl.king.im/king-lib/framework/internal/stat"
	fwTracing "gl.king.im/king-lib/framework/internal/tracing"
	mauth "gl.king.im/king-lib/framework/middlewares/auth"
	"gl.king.im/king-lib/framework/service"
	"gl.king.im/king-lib/framework/service/desc"
	fwMd "gl.king.im/king-lib/framework/transport/metadata"
	admPb "gl.king.im/king-lib/protobuf/api/admin/service/v1"
	"gl.king.im/king-lib/protobuf/api/common"
	"google.golang.org/protobuf/proto"

	kgin "github.com/go-kratos/gin"
	ke "github.com/go-kratos/kratos/v2/errors"
)

func ginErrHandler(ctx *gin.Context, err error) {
	//fmt.Println("ginErrHandler::", err)
	if err != nil {
		ret := map[string]interface{}{}
		httpStatus := 400
		switch typeErr := err.(type) {
		case *ke.Error:
			if typeErr.Code < 1000 {
				httpStatus = int(typeErr.Code)
			}

			ret = map[string]interface{}{
				"code":     typeErr.Code,
				"reason":   typeErr.Reason,
				"message":  typeErr.Error(),
				"metadata": typeErr.Metadata,
			}
		default:
			ret = map[string]interface{}{
				"code":    400,
				"message": typeErr.Error(),
			}
		}
		ctx.AbortWithStatusJSON(httpStatus, ret)
		//fmt.Println("after AbortWithStatusJSON::")
		//ctx.JSON(400, ret)
	}
	//fmt.Println("ginErrHandler end")
}

func errHandler(frameworkType string, ginCtx *gin.Context, err error) (interface{}, error) {
	if frameworkType == "gin" {
		ginErrHandler(ginCtx, err)
	}
	return nil, err
}

type FrameworkMwCnf struct {
	Type                    string
	ContextDeadlineDuration time.Duration
}

var serviceInExternalNetwork bool = false

func SetServiceExternalNetworkState(yes bool) {
	serviceInExternalNetwork = yes
}

var userAuthMethod config.Service_AuthMethod_Enums
var svcAuthMethod config.Service_AuthMethod_Enums

func SetUserAuthMethod(am config.Service_AuthMethod_Enums) {
	userAuthMethod = am
}

func SetSvcAuthMethod(am config.Service_AuthMethod_Enums) {
	svcAuthMethod = am
}

func BaseMiddleWareCore(next middleware.Handler, fMwCnf FrameworkMwCnf) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, replyErr error) {
		logFlag := "AuthMiddleWareCore:"
		log.Println(logFlag)

		ctx = stat.NewCtxStatServerContext(ctx, stat.NewCtxStat())
		ctx = fwTracing.NewGlobalTraceInfoServerContext(ctx, &fwTracing.GlobalTraceInfo{})

		var ginCtx *gin.Context

		frameworkType := fMwCnf.Type

		if frameworkType == "gin" {
			var ok bool
			ginCtx, ok = kgin.FromGinContext(ctx)

			if fMwCnf.ContextDeadlineDuration != 0 {
				//var cancelFun context.CancelFunc
				//d1, _ := ctx.Deadline()
				//ctx, cancelFun = context.WithTimeout(ctx, fMwCnf.ContextDeadlineDuration)
				//_ = cancelFun
				//d2, _ := ctx.Deadline()
				//fmt.Println(logFlag, ";fMwCnf.ContextDeadlineDuration:", fMwCnf.ContextDeadlineDuration, ";d1::", d1, ";d2::", d2)
			}

			_ = ok
		}

		//trace id
		traceIdValuer := tracing.TraceID()
		traceId := traceIdValuer(ctx).(string)
		log.Println(logFlag, "traceId::", traceId)
		//ctx = metadata.AppendToClientContext(ctx, "x-md-global-trace-id", traceId)

		//transport断言
		tr, ok := transport.FromServerContext(ctx)
		if !ok {
			replyErr = ke.New(400, "BASE_MW_TRANSPORT_ERROR", "TRANSPORT CONTEXT ERROR")
			return errHandler(frameworkType, ginCtx, replyErr)
		}

		trKind := tr.Kind()
		var httpTr *khttp.Transport
		var httpReq *http.Request

		if trKind == transport.KindHTTP {
			httpTr = tr.(*khttp.Transport)
			httpReq = httpTr.Request()

			httpTr.ReplyHeader().Set("X-TRACE-ID", traceId)
		}
		header := tr.RequestHeader()
		log.Println("=========framework BaseMiddleWare Header..=============.", header)

		//被调用接口名
		op := tr.Operation()
		log.Println(logFlag, "op:", op, "tr kind:", trKind)

		if op == "" {
			kerr := ke.New(400, "BASE_MW_TRANSPORT_OPERATION_EMPTY_ERROR", "TRANSPORT OPERATION EMPTY")
			return errHandler(frameworkType, ginCtx, kerr)
		}

		pbBossOpts := pBAuthOptsParse(op)
		//目前来看认证忽略就等同于认证+鉴权都忽略，因为认证都忽略了，意味着用户实体是可能不存在，而用户实体不存在那鉴权也就没有意义了
		//也就是说这里可以只需要存在(authenticationIgnore或authIgnore)和authorizationIgnore，因为authenticationIgnore等同于authIgnore，这里推荐使用authenticationIgnore
		authIgnore, authenticationIgnore, authorizationIgnore := getAuthIgnores(pbBossOpts)

		serviceConf := config.GetServiceConf()
		//注意这里不能直接authIgnore = serviceConf.Service.AuthClose，因为比如authIgnore原本是true，这里直接赋值就会被覆盖为false，是不符合预期逻辑的
		if serviceConf.Service.AuthClose {
			authIgnore = true
		}

		md, ok := metadata.FromServerContext(ctx)
		log.Println(logFlag, "md, ok:", md, ok)
		if !ok {
			log.Println("AuthMiddleWareCore BASE_MW_CALLER_INFO_ERROR:", "元信息获取失败！")
			replyErr = ke.New(400, "BASE_MW_CALLER_INFO_ERROR", "元信息获取失败！")
			return errHandler(frameworkType, ginCtx, replyErr)
		}

		reqAuthMethods := auth.ParseAuthMethod(ctx)

		var firstAuthMethod *auth.AuthMethod
		if len(reqAuthMethods) > 0 {
			firstAuthMethod = &reqAuthMethods[0]
		}

		//这里注意，当authIgnore为true时，就不做鉴权相关的报错了
		//连header/metadata的检查也不能做了，因为有很多接口是供外部系统调用的，不一定遵循我们定义的header/metadata规则
		//同时注意，如果authIgnore为true则user实体也不会进行获取和设置
		if !authIgnore {
			if len(reqAuthMethods) < 1 {
				return errHandler(frameworkType, ginCtx, ke.Unauthorized("AUTH_METHOD_EMPTY_ERROR", "未认证"))
			}

			if len(reqAuthMethods) > 1 {
				return errHandler(frameworkType, ginCtx, ke.BadRequest("AUTH_METHOD_ERROR", "不支持同时使用多种认证方式"))
			}

			err, user := authCheckForBaseMw(ctx, op, firstAuthMethod, frameworkType, pbBossOpts, authIgnore, authenticationIgnore, authorizationIgnore)
			if err != nil {
				return errHandler(frameworkType, ginCtx, err)
			}

			ctx = authUser.NewUserServerContext(ctx, user)

			//local md做鉴权使用、global md做信息传递使用，不能用global作为鉴权，否则所有服务都会以global的token做鉴权，也就是前端传递也只能传local不能传global
			switch strings.ToLower(firstAuthMethod.MethodKey) {
			case strings.ToLower(framework.METADATA_KEY_AUTH_TOKEN):
				switch strings.ToLower(firstAuthMethod.SubMethod) {
				case strings.ToLower(framework.MDV_SERVICE_CALL_SCENARIO_INNER): //服务内部场景

				case strings.ToLower(framework.MDV_SERVICE_CALL_SCENARIO_GETEWAY): //网关场景
					ctx = writeOpUserMD(ctx, firstAuthMethod.Token)
				}
			//取代原来的gateway
			case strings.ToLower(framework.METADATA_KEY_LOCAL_OP_USER_TOKEN):
				ctx = writeOpUserMD(ctx, firstAuthMethod.Token)

			//取代原来的inner
			case strings.ToLower(framework.METADATA_KEY_LOCAL_SVC_TOKEN):

			}

		}

		if firstAuthMethod != nil {
			switch strings.ToLower(firstAuthMethod.MethodKey) {
			case strings.ToLower(framework.METADATA_KEY_AUTH_TOKEN):
				switch strings.ToLower(firstAuthMethod.SubMethod) {
				case strings.ToLower(framework.MDV_SERVICE_CALL_SCENARIO_INNER): //服务内部场景

				case strings.ToLower(framework.MDV_SERVICE_CALL_SCENARIO_GETEWAY): //网关场景
					ctx = fwTracing.NewGlobalTraceInfoServerContext(ctx, &fwTracing.GlobalTraceInfo{IsUserBeginNode: true})
				}
			//取代原来的gateway
			case strings.ToLower(framework.METADATA_KEY_LOCAL_OP_USER_TOKEN):
				ctx = fwTracing.NewGlobalTraceInfoServerContext(ctx, &fwTracing.GlobalTraceInfo{IsUserBeginNode: true})

			//取代原来的inner
			case strings.ToLower(framework.METADATA_KEY_LOCAL_SVC_TOKEN):

			}
		}

		//设置ocean token服务间使用
		globalInfo, ok := authUser.OpUserGlobalInfoFromServerContext(ctx)
		if ok {
			ctx = authUser.NewOceanTokenCtx(ctx, globalInfo.OceanAuthToken)
		}

		//目前只对boss登录，即gateway scenario记admin log

		if userAuthMethod == config.Service_AuthMethod_Default && len(reqAuthMethods) == 1 && reqAuthMethods[0].MethodKey == framework.METADATA_KEY_AUTH_TOKEN && reqAuthMethods[0].SubMethod == framework.MDV_SERVICE_CALL_SCENARIO_GETEWAY {
			admLog(ctx, req, httpReq, op, frameworkType, md)
		}

		if replyErr == nil {
			rsp, err := next(ctx, req)

			return rsp, err
		} else {
			if frameworkType == "gin" {
				ginErrHandler(ginCtx, replyErr)
			}
			return nil, replyErr
		}
	}
}

// 写入op用户的全局传递token以及一些运营用户全局信息
func writeOpUserMD(ctx context.Context, token string) context.Context {

	ctx = metadata.AppendToClientContext(ctx, strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_TOKEN), token)
	opUser, ok := authUser.BossOpUserFromServerContext(ctx)
	if ok {
		ctx = metadata.AppendToClientContext(ctx, strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_USERNAME), opUser.Username)
		//由于公司k8s集群内会报： header key \"x-md-global-op-user-name\" contains value with non-printable ASCII characters 的错误，暂时注释掉姓名字段的传递
		//ctx = metadata.AppendToClientContext(ctx, strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_NAME), opUser.Name)
		ctx = metadata.AppendToClientContext(ctx, strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_ID), fmt.Sprintf("%d", opUser.Id))
		ctx = metadata.AppendToClientContext(ctx, strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_OCEAN_TOKEN), opUser.OceanAuthToken)

		if opUser.SessionInfo != nil {
			ctx = metadata.AppendToClientContext(ctx, strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_LOGIN_METHOD), opUser.SessionInfo.LoginMethod.String())
			ctx = metadata.AppendToClientContext(ctx, strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_CONSOLE_TOKEN), opUser.SessionInfo.ConsoleAuthToken)
			ctx = metadata.AppendToClientContext(ctx, strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_FEISHU_TOKEN), opUser.SessionInfo.FeishuAuthToken)
		}
	}
	return ctx
}

// 鉴权方式处理
func authCheckForBaseMw(ctx context.Context, op string, validAuthMethod *auth.AuthMethod, frameworkType string, pbBossOpts *common.BossOpts, authIgnore, authenticationIgnore, authorizationIgnore bool) (err error, user interface{}) {
	logFlag := "authCheckForBaseMw:"
	switch validAuthMethod.MethodKey {
	case framework.METADATA_KEY_AUTH_TOKEN:
		if validAuthMethod.SubMethod == "" {
			log.Println(logFlag, "BASE_MW_CALLER_INFO_ERROR:", "调用者信息缺失！")
			err = ke.New(401, "BASE_MW_CALLER_INFO_ERROR", "调用者信息缺失！")
			return
		}

		log.Println(logFlag, "X-Md-Global-Call-Scenario:", validAuthMethod.SubMethod)

		//X-Md-Global-Call-Scenario的检查，注意这个检查只能放在framework.METADATA_KEY_AUTH_TOKEN分支里面，因为别的场景可能不按我们定义的header来请求，比如收银台、电子合同的回调请求
		if frameworkType == "kratos" {
			err = checkMethodOpen(pbBossOpts, validAuthMethod.SubMethod)
			if err != nil {
				return
			}
		}

		//注意这里不能直接将函数返回结果赋值给authorizationIgnore，比如authorizationIgnore原本是true，而这里返回结果是false，就会误覆盖authorizationIgnore原本的值
		hasSet := mauth.HasSetAuthorizationIgnore(op)

		if hasSet {
			authorizationIgnore = true
		}

		//session.SetGlobalAuthToken(authToken)
		service.AppInfoIns.CallScenario = validAuthMethod.SubMethod
		service.AppInfoIns.Caller = validAuthMethod.SubMethod

		switch validAuthMethod.SubMethod {
		case framework.MDV_SERVICE_CALL_SCENARIO_INNER: //服务内部场景
			if svcAuthMethod == config.Service_AuthMethod_OutsideHttp {
				err, user = mauth.InnerAuthCheckInExternalNetwork(ctx, validAuthMethod.Token, op, authIgnore, authenticationIgnore, authorizationIgnore)
			} else {
				err, user = mauth.AuthTokenInner(ctx, validAuthMethod.Token, op, authIgnore, authenticationIgnore, authorizationIgnore)
			}

			if err != nil {
				return
			}

			if user != nil {
				assertUser, assertOk := user.(*common.BossService)
				if assertOk {
					log.Println(logFlag, "from service:", assertUser.Name, assertUser.Platform.String())
				}
			}

		case framework.MDV_SERVICE_CALL_SCENARIO_GETEWAY: //网关场景
			if userAuthMethod == config.Service_AuthMethod_OutsideHttp {
				err, user = mauth.GatewayScenarioAuthCheckInExternalNetwork(ctx, validAuthMethod.Token, op, authIgnore, authenticationIgnore, authorizationIgnore)
			} else {
				err, user = mauth.AuthTokenGatewayScenario(ctx, validAuthMethod.Token, op, authIgnore, authenticationIgnore, authorizationIgnore)
			}

			if err != nil {
				return
			}

			if user != nil {
				assertUser, assertOk := user.(*common.BossOperationUser)
				_ = assertUser
				if assertOk {
					//log.Println(logFlag, "from user:", assertUser.Username, assertUser.Name)
				}
			}

		default:
			err, user = mauth.AuthTokenDefault(ctx, validAuthMethod.Token, op, authIgnore, authenticationIgnore, authorizationIgnore, validAuthMethod.SubMethod)
			if err != nil {
				return
			}
		}
	case framework.METADATA_KEY_PRISM_ACCESS_TOKEN:

		err, user = prismAccessToken(ctx, op, authIgnore, authenticationIgnore, authorizationIgnore, validAuthMethod.Token)
		if err != nil {
			//return nil, err
			return
		}
	case framework.METADATA_KEY_OPEN_ACCESS_TOKEN:
		err, user = mauth.OpenapiAuth(ctx, validAuthMethod.Token, op, authIgnore, authenticationIgnore, authorizationIgnore)
		if err != nil {
			return
		}

		//err, user = openapiAuth(ctx, validAuthMethod.MethodKey, op, authIgnore, authenticationIgnore, authorizationIgnore, validAuthMethod.Token)
		//if err != nil {
		//	//return nil, err
		//	return
		//}
	//取代原来的gateway
	case framework.METADATA_KEY_LOCAL_OP_USER_TOKEN:
		if userAuthMethod == config.Service_AuthMethod_OutsideHttp {
			err, user = mauth.GatewayScenarioAuthCheckInExternalNetwork(ctx, validAuthMethod.Token, op, authIgnore, authenticationIgnore, authorizationIgnore)
		} else {
			err, user = mauth.AuthTokenGatewayScenario(ctx, validAuthMethod.Token, op, authIgnore, authenticationIgnore, authorizationIgnore)
		}

		if err != nil {
			return
		}

		if user != nil {
			assertUser, assertOk := user.(*common.BossOperationUser)
			_ = assertUser
			if assertOk {
				//log.Println(logFlag, "from user:", assertUser.Username, assertUser.Name)
			}
		}

	//取代原来的inner
	case framework.METADATA_KEY_LOCAL_SVC_TOKEN:
		if svcAuthMethod == config.Service_AuthMethod_OutsideHttp {
			err, user = mauth.InnerAuthCheckInExternalNetwork(ctx, validAuthMethod.Token, op, authIgnore, authenticationIgnore, authorizationIgnore)
		} else {
			err, user = mauth.AuthTokenInner(ctx, validAuthMethod.Token, op, authIgnore, authenticationIgnore, authorizationIgnore)
		}

		if err != nil {
			return
		}

		if user != nil {
			assertUser, assertOk := user.(*common.BossService)
			if assertOk {
				log.Println(logFlag, "from service:", assertUser.Name, assertUser.Platform.String())
			}
		}

	default:
		return ke.BadRequest("NOT_SUPPORTED_AUTH_METHOD", "不支持的认证方式"), user
	}

	return
}

func pBAuthOptsParse(op string) (pbBossOpts *common.BossOpts) {
	logFlag := "PBAuthOptsParse::"
	opSplit := strings.Split(op, "/")
	methodName := opSplit[len(opSplit)-1]
	serviceName := strings.TrimRight(op, methodName)
	serviceName = strings.Trim(serviceName, "/")
	srvRegConf, srvRegCnfOk := desc.ServerRegCnfMap[serviceName]

	log.Println(logFlag, "operation info:", op, opSplit, serviceName, methodName)

	//获取接口的option定义
	if srvRegCnfOk {
		methodsDesc := srvRegConf.MethodMap[methodName]
		pbBossOpts = proto.GetExtension(methodsDesc.Options, common.E_BossOpts).(*common.BossOpts)
	}

	return
}

// 判断接口的开放域
func checkMethodOpen(pbBossOpts *common.BossOpts, caller string) error {

	if pbBossOpts == nil {

	} else {
		var in bool
		//如果不填写则都开放
		if len(pbBossOpts.MethodOpen) == 0 {
			in = true
		}

		for _, v := range pbBossOpts.MethodOpen {
			if caller == v {
				in = true
				break
			}
		}

		if !in {
			kerr := ke.New(400, "BASE_MW_OPERATION_SCENARIO_ERROR", "该接口不适用此调用场景")
			return kerr
		}

	}

	return nil
}

func getAuthIgnores(pbBossOpts *common.BossOpts) (authIgnore bool, authenticationIgnore bool, authorizationIgnore bool) {
	//获取接口的option定义
	if pbBossOpts == nil {

	} else {

		//是否忽略认证和权限校验-如果是则直接跳出
		authIgnore = pbBossOpts.AuthIgnore

		//是否忽略认证校验
		authenticationIgnore = pbBossOpts.AuthenticationIgnore

		//是否忽略权限校验
		authorizationIgnore = pbBossOpts.AuthorizationIgnore
	}

	return
}

func prismAccessToken(ctx context.Context, op string, authIgnore, authenticationIgnore, authorizationIgnore bool, prismAccessToken string) (err error, user interface{}) {
	if authIgnore {
		return
	}

	val, ok := service.AppInfoIns.AuthMethodRegistry.Load(framework.METADATA_KEY_PRISM_ACCESS_TOKEN)
	if ok {

		ar := val.(auth.AuthRegister)
		//认证
		if !authenticationIgnore {
			var validateOk bool
			validateOk, user, err = ar.Ae.Validate(ctx, prismAccessToken)
			if err != nil {
				return
			}

			if !validateOk {
				err = ke.New(401, "BASE_MW_AUTHENTICATION_ERROR", "当前未认证")
				return
			}
		}
		//权限
		if !authorizationIgnore {
			var aoOk bool
			aoOk, err = ar.Ao.Can(ctx, prismAccessToken, op)
			if err != nil {
				return
			}

			if !aoOk {
				err = ke.New(403, "BASE_MW_AUTHORIZATION_ERROR", "此操作未授权")
				return
			}
		}
	} else {
		err = ke.New(400, "BASE_MW_CALLER_INFO_ERROR", "不支持的服务开放域！")
		return
	}

	return
}

func openapiAuth(ctx context.Context, authMethodKey string, op string, authIgnore, authenticationIgnore, authorizationIgnore bool, token string) (err error, user interface{}) {
	if authIgnore {
		return
	}

	val, ok := service.AppInfoIns.AuthMethodRegistry.Load(authMethodKey)
	if ok {

		ar := val.(auth.AuthRegister)
		//认证
		if !authenticationIgnore {
			var validateOk bool
			validateOk, user, err = ar.Ae.Validate(ctx, token)
			if err != nil {
				return
			}

			if !validateOk {
				err = ke.New(401, "BASE_MW_AUTHENTICATION_ERROR", "当前未认证")
				return
			}
		}
		//权限
		if !authorizationIgnore {
			var aoOk bool
			aoOk, err = ar.Ao.Can(ctx, token, op)
			if err != nil {
				return
			}

			if !aoOk {
				err = ke.New(403, "BASE_MW_AUTHORIZATION_ERROR", "此操作未授权")
				return
			}
		}
	} else {
		err = ke.New(400, "UnspportedAuthMethod", "不支持的认证方式")
		return
	}

	return
}

func admLog(ctx context.Context, req interface{}, httpReq *http.Request, op string, frameworkType string, md metadata.Metadata) {
	logFlag := "middleware AdmLog::"
	admLogMsg := ""
	reqBytes, _ := json.Marshal(req)
	admLogMsg += string(reqBytes)
	admLog := &admPb.AdminLog{
		//Method: ,
		Operation: op,
		Message:   admLogMsg,
		Metadata:  fmt.Sprintf("%+v", md),
	}

	if httpReq != nil {
		admLog.Method = httpReq.Method
		admLog.Path = httpReq.URL.Path
		admLog.Params = httpReq.URL.RawQuery
		admLog.Host = httpReq.Host
		admLog.HttpRaw = fmt.Sprintf("%+v", httpReq)
		admLog.Framework = frameworkType
		//if frameworkType == "gin" {
		//	httpReq.ParseForm()
		//	admLog.Message = fmt.Sprintf("%+v", httpReq.PostForm)
		//}
	}

	userIdentity, ok := authUser.BossOpUserFromServerContext(ctx)
	if ok && userIdentity != nil {
		admLog.Username = userIdentity.Username
	}

	go func(al *admPb.AdminLog) {
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				log.Println("AdminLog Recovery :", recoverErr)
			}
		}()

		admCli, err := data.GetAdmServiceClient()
		if err != nil {
			log.Println(logFlag, "data.GetAdmServiceClient() err:", err)
			return
		}
		callCtx := fwMd.BuildInnerMDCtx()

		_, err = admCli.AdminLog(callCtx, &admPb.AdminLogRequest{
			Logs: []*admPb.AdminLog{
				al,
			},
		})

		if err != nil {
			log.Println(logFlag, "admCli.AdminLog() err:", err)
		}
	}(admLog)
}
