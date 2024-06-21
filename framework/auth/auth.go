package auth

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/metadata"
	"gl.king.im/king-lib/framework"
)

type AuthRegister struct {
	Ae Authentication
	Ao Authorization
}

// 认证+授权接口
type Auth interface {
	//GetToken() (token string,err error)        //作为调用方时生成/获取自身的token，以用于发起调用给被调用方鉴权
	AuthCheck(ctx context.Context, token, resource string, aeIgnore, aoIgnore bool) (aeOk, aoOk bool, user interface{}, err error) //作为被调用方的认证校验
}

// 认证接口
type Authentication interface {
	//GetToken() (token string,err error)        //作为调用方时生成/获取自身的token，以用于发起调用给被调用方鉴权
	Validate(ctx context.Context, token string) (ok bool, user interface{}, err error) //作为被调用方的认证校验
}

// 授权接口
type Authorization interface {
	Can(ctx context.Context, token, resource string) (bool, error) //作为被调用方的授权校验
}

// 认证接口空实现
var _ Authentication = (*UselessAuthentication)(nil)

type UselessAuthentication struct {
}

func (*UselessAuthentication) GetToken() (token string) {
	token = "UselessAuthenticationToken"
	return token
}

func (*UselessAuthentication) Validate(ctx context.Context, token string) (ok bool, user interface{}, err error) {
	ok = true
	return ok, user, err
}

// 授权接口空实现
var _ Authorization = (*UselessAuthorization)(nil)

type UselessAuthorization struct {
}

func (*UselessAuthorization) Can(ctx context.Context, token, resource string) (bool, error) {
	return true, nil
}

type AuthMethod struct {
	MethodKey string
	Token     string
	SubMethod string
}

func ParseAuthMethod(ctx context.Context) (reqAuthMethods []AuthMethod) {

	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		return
	}

	prismAccessToken, hasPrismATMethod := md[strings.ToLower(framework.METADATA_KEY_PRISM_ACCESS_TOKEN)]
	if hasPrismATMethod {
		reqAuthMethods = append(reqAuthMethods, AuthMethod{
			MethodKey: framework.METADATA_KEY_PRISM_ACCESS_TOKEN,
			Token:     prismAccessToken,
		})
		//validAuthMethod = framework.METADATA_KEY_PRISM_ACCESS_TOKEN
	}

	globalAuthToken, hasAuthTokenMethod := md[strings.ToLower(framework.METADATA_KEY_AUTH_TOKEN)]
	if hasAuthTokenMethod {
		caller := md.Get(framework.METADATA_KEY_CALL_SCENARIO)
		reqAuthMethods = append(reqAuthMethods, AuthMethod{
			MethodKey: framework.METADATA_KEY_AUTH_TOKEN,
			Token:     globalAuthToken,
			SubMethod: caller,
		})
		//validAuthMethod = framework.METADATA_KEY_AUTH_TOKEN
	}

	openapiAccessToken, hasOpenapiATMethod := md[strings.ToLower(framework.METADATA_KEY_OPEN_ACCESS_TOKEN)]
	if hasOpenapiATMethod {
		reqAuthMethods = append(reqAuthMethods, AuthMethod{
			MethodKey: framework.METADATA_KEY_OPEN_ACCESS_TOKEN,
			Token:     openapiAccessToken,
		})
	}

	//op user的local md做鉴权使用、global md做信息传递使用，不能用global作为鉴权，否则所有服务都会以global的token做鉴权，也就是前端传递也只能传local不能传global
	opUserAccessToken, hasOpUserAuthTokenMethod := md[strings.ToLower(framework.METADATA_KEY_LOCAL_OP_USER_TOKEN)]
	if hasOpUserAuthTokenMethod {
		reqAuthMethods = append(reqAuthMethods, AuthMethod{
			MethodKey: framework.METADATA_KEY_LOCAL_OP_USER_TOKEN,
			Token:     opUserAccessToken,
		})
	}

	svcAccessToken, hasSvcAuthTokenMethod := md[strings.ToLower(framework.METADATA_KEY_LOCAL_SVC_TOKEN)]
	if hasSvcAuthTokenMethod {
		reqAuthMethods = append(reqAuthMethods, AuthMethod{
			MethodKey: framework.METADATA_KEY_LOCAL_SVC_TOKEN,
			Token:     svcAccessToken,
		})
	}

	return
}
