package user

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	ke "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"gl.king.im/king-lib/framework"
	"gl.king.im/king-lib/framework/auth"
	"gl.king.im/king-lib/framework/auth/token"
	"gl.king.im/king-lib/framework/internal/data"
	"gl.king.im/king-lib/framework/internal/tracing"
	"gl.king.im/king-lib/framework/session"
	"gl.king.im/king-lib/framework/session/user"
	fwMd "gl.king.im/king-lib/framework/transport/metadata"
	cpb "gl.king.im/king-lib/protobuf/api/common"
	servicePb "gl.king.im/king-lib/protobuf/api/service/service/v1"
)

// 尽可能的囊括所有认证方式的用户实体结构--计划废弃
type UserIdentity struct {
	Id            uint64
	Username      string
	Secret        string
	Platform      uint32 //业务平台，0为未知，1为BOSS，2为控制台
	Email         string
	Phone         string
	Password      string
	LastLoginTime string
	LdapDn        string
	LdapCn        string
}

func GetIdentityFromContext(ctx context.Context) (*UserIdentity, error) {

	scenario, err := session.GetContextScenario(ctx)
	if err != nil {
		return nil, ke.New(400, "Scenario_FAIL", err.Error())
	}

	authToken, err := token.GetTokenFromContext(ctx)
	if err != nil {
		return nil, ke.New(400, "SessionId_FAIL", err.Error())
	}

	switch scenario {
	case framework.MDV_SERVICE_CALL_SCENARIO_GETEWAY: //boss默认网关，走的global sesssion

		ssUserIdentity, err := user.GetUserIdentity(ctx)
		if err != nil {
			return nil, err
		}

		//identity := &session.UserIdentity{}
		//err := session.GetIdentityByToken(authToken, identity)
		//if err != nil {
		//	return nil, err
		//}

		user := &UserIdentity{
			Username: ssUserIdentity.Username,
			Platform: 1,
		}

		return user, nil
	case framework.MDV_SERVICE_CALL_SCENARIO_INNER: //内部互访，走的jwt
		client, err := data.GetServiceServiceClient()
		if err != nil {
			return nil, err
		}
		callCtx := fwMd.BuildInnerMDCtx()

		//connGRPC, callCtx, err := service.NewGrpcClientConn("BossService")
		//if err != nil {
		//	panic(err)
		//}
		//defer connGRPC.Close()
		//client := servicePb.NewServiceServiceClient(connGRPC)

		serviceInfo, err := client.GetServiceIdentityByToken(callCtx, &servicePb.GetServiceIdentityByTokenRequest{Jwt: authToken})

		if err == nil {
			user := &UserIdentity{
				Username: serviceInfo.User.Name,
				Platform: serviceInfo.User.Platform,
			}
			return user, nil
		} else {
			return nil, err
		}
	default:
		return nil, ke.New(400, "Scenario_FAIL", "不支持的调用场景！")
	}
}

type userCtxKey struct{}

// NewServerContext creates a new context with client md attached.
func NewUserServerContext(ctx context.Context, user interface{}) context.Context {
	return context.WithValue(ctx, userCtxKey{}, user)
}

// FromServerContext returns the server metadata in ctx if it exists.
func UserFromServerContext(ctx context.Context) interface{} {
	user := ctx.Value(userCtxKey{})
	return user
}

type bossOpuserCtxKey struct{}

type bossServiceUserCtxKey struct{}

func BossOpUserFromServerContext(ctx context.Context) (user *cpb.BossOperationUser, ok bool) {
	user, ok = ctx.Value(userCtxKey{}).(*cpb.BossOperationUser)
	return user, ok
}

func ServiceUserFromServerContext(ctx context.Context) (user *cpb.BossService, ok bool) {
	user, ok = ctx.Value(userCtxKey{}).(*cpb.BossService)
	return user, ok
}

func OpenapiUserFromServerContext(ctx context.Context) (user *cpb.BossService, ok bool) {
	user, ok = ctx.Value(userCtxKey{}).(*cpb.BossService)
	return user, ok
}

func UsernameLogValuer() log.Valuer {
	return func(ctx context.Context) interface{} {
		reqAuthMethods := auth.ParseAuthMethod(ctx)

		if len(reqAuthMethods) > 0 && reqAuthMethods[0].MethodKey == framework.METADATA_KEY_AUTH_TOKEN {
			var username string
			switch reqAuthMethods[0].SubMethod {
			case framework.MDV_SERVICE_CALL_SCENARIO_GETEWAY:
				userEntity, ok := BossOpUserFromServerContext(ctx)

				if ok {
					username = userEntity.Username
				}
			case framework.MDV_SERVICE_CALL_SCENARIO_INNER:
				userEntity, ok := ServiceUserFromServerContext(ctx)
				if ok {
					username = userEntity.Name
				}

			}

			return username
		}
		return ""
	}
}

//由于调用方和被调用方获取metadata的方案不相同，这里需要做判断

type OpUserGlobalInfo struct {
	AuthToken      string
	Username       string
	Name           string
	UserId         uint64
	OceanAuthToken string
	OceanUserId    int64
	SessionInfo    cpb.SessionInfo
}

func OpUserGlobalInfoFromServerContext(ctx context.Context) (info *OpUserGlobalInfo, ok bool) {

	globalTraceInfo, ok := tracing.GlobalTraceInfoFromServerContext(ctx)
	if !ok {
		return
	}

	var md metadata.Metadata
	if globalTraceInfo.IsUserBeginNode {
		md, ok = metadata.FromClientContext(ctx)
		if !ok {
			fmt.Println("OpUserGlobalInfoFromServerContext:", "元信息获取失败！")
			return
		}
	} else {
		md, ok = metadata.FromServerContext(ctx)
		if !ok {
			fmt.Println("OpUserGlobalInfoFromServerContext:", "元信息获取失败！")
			return
		}
	}

	info = &OpUserGlobalInfo{}

	if val, mdOk := md[strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_TOKEN)]; mdOk {
		info.AuthToken = val
	}

	if val, mdOk := md[strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_OCEAN_TOKEN)]; mdOk {
		info.OceanAuthToken = val
	}

	if val, mdOk := md[strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_USERNAME)]; mdOk {
		info.Username = val
	}

	if val, mdOk := md[strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_NAME)]; mdOk {
		info.Name = val
	}

	if val, mdOk := md[strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_ID)]; mdOk {
		uidInt, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			fmt.Println("OpUserGlobalInfoFromServerContext strconv.ParseUint err", err)
		} else {
			info.UserId = uidInt
		}
	}

	if val, mdOk := md[strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_LOGIN_METHOD)]; mdOk {
		if lgVal, lgOk := cpb.LoginMethod_Enums_value[val]; lgOk {
			info.SessionInfo.LoginMethod = cpb.LoginMethod_Enums(lgVal)
		}
	}

	if val, mdOk := md[strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_FEISHU_TOKEN)]; mdOk {
		info.SessionInfo.FeishuAuthToken = val
	}

	if val, mdOk := md[strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_CONSOLE_TOKEN)]; mdOk {
		info.SessionInfo.ConsoleAuthToken = val
	}

	if val, mdOk := md[strings.ToLower(framework.METADATA_KEY_GLOBAL_OP_USER_OCEAN_TOKEN)]; mdOk {
		info.SessionInfo.OceanAuthToken = val
	}

	return
}

type oceanTokenCtxKey struct{}

const dpPrefix = "data_platform:"

// NewOceanTokenCtx creates a new context with client md attached.
func NewOceanTokenCtx(ctx context.Context, user interface{}) context.Context {
	return context.WithValue(ctx, oceanTokenCtxKey{}, user)
}

func GlobalOceanTokenFromServerContext(ctx context.Context) (token string, ok bool) {
	token, ok = ctx.Value(oceanTokenCtxKey{}).(string)
	return dpPrefix + token, ok
}
