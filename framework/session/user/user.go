package user

import (

	//mcrRdsStore "github.com/go-macaron/session/redis"

	"bytes"
	"context"
	"encoding/gob"
	"log"

	_ "github.com/astaxie/beego/session/redis"
	ke "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metadata"
	"gl.king.im/king-lib/framework"
	"gl.king.im/king-lib/framework/internal/data"
	fwMd "gl.king.im/king-lib/framework/transport/metadata"
	userPb "gl.king.im/king-lib/protobuf/api/user/service/v1"
)

type UserIdentity struct {
	Id       uint64
	Username string
	Nickname string
	Avatar   string
	Secret   string
	Platform uint32 //业务平台，0为未知，1位OA/LDAP，2为BOSS
	Email    string
	Phone    string
	//Password      string
	LastLoginTime string
	LdapDn        string
	LdapCn        string
}

func getValidTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromServerContext(ctx)
	log.Println("GetContextSessionId...md, ok:", md, ok)
	if !ok {
		return "", ke.New(400, "METADATA_INFO_ERROR", "元信息获取失败！")
	}

	authToken := md.Get(framework.METADATA_KEY_AUTH_TOKEN)
	if authToken == "" {
		return "", ke.New(400, "METADATA_AUTH_TOKEN_EMPTY", "用户session_id为空！")
	}

	return authToken, nil
}

type userIdentityOptions struct {
	dbSource bool
}
type userIdentityOption func(*userIdentityOptions)

func WithDbSource(dbsrc bool) userIdentityOption {
	return func(o *userIdentityOptions) {
		o.dbSource = dbsrc
	}
}

func GetUserIdentity(reqCtx context.Context, opts ...userIdentityOption) (*UserIdentity, error) {
	authToken, err := getValidTokenFromContext(reqCtx)
	if err != nil {
		return nil, err
	}

	return GetUserIdentityByToken(reqCtx, authToken, opts...)
}

func GetUserIdentityByToken(reqCtx context.Context, authToken string, opts ...userIdentityOption) (*UserIdentity, error) {

	client, err := data.GetUserServiceClient()
	if err != nil {
		return nil, err
	}
	//callCtx := fwMd.BuildInnerMDCtx()

	//	connGRPC, callCtx, err := service.NewGrpcClientConn("BossUser")
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	defer connGRPC.Close()
	//
	//	client := userPb.NewUserServiceClient(connGRPC)

	options := &userIdentityOptions{}
	for _, o := range opts {
		o(options)
	}

	//数据库实时数据
	//dbSource := false
	//if len(args) > 0 {
	//	dbSource = args[0].(bool)
	//}

	var replyIdentity *userPb.SessionUserIdentity
	if options.dbSource {
		reply, err := client.GetUserIdentityBySessionId(reqCtx, &userPb.GetUserIdentityBySessionIdRequest{SessionId: authToken})
		if err != nil {
			return nil, err
		}

		replyIdentity = reply.Identity
	} else {
		reply, err := client.GetSessionUserIdentity(reqCtx, &userPb.GetSessionUserIdentityRequest{SessionId: authToken})
		if err != nil {
			return nil, err
		}
		replyIdentity = reply.Identity
	}

	if replyIdentity == nil {
		return nil, nil
	}

	identity := &UserIdentity{
		Id:       replyIdentity.Id,
		Username: replyIdentity.Username,
		Secret:   replyIdentity.Secret,
		Nickname: replyIdentity.Nickname,
		Avatar:   replyIdentity.Avatar,
		//Password:      replyIdentity.Password,
		Email:         replyIdentity.Email,
		Phone:         replyIdentity.Phone,
		LastLoginTime: replyIdentity.LastLoginTime,
	}

	return identity, nil
}

func SessionGet(ctx context.Context, key string, scanValue interface{}) (interface{}, error) {
	token, err := getValidTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	client, err := data.GetUserServiceClient()
	if err != nil {
		return nil, err
	}
	callCtx := fwMd.BuildInnerMDCtx()

	//	connGRPC, callCtx, err := service.NewGrpcClientConn("BossUser")
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	defer connGRPC.Close()
	//
	//	client := userPb.NewUserServiceClient(connGRPC)

	reply, err := client.SessionGet(callCtx, &userPb.SessionGetRequest{SessionId: token, Key: key})
	if err != nil {
		return nil, err
	}

	if reply.Value == nil {
		return nil, nil
	}

	valueBuf := bytes.NewBuffer(reply.Value)
	dec := gob.NewDecoder(valueBuf)

	err = dec.Decode(scanValue)

	if err != nil {
		return nil, err
	}

	return reply.Value, nil
}

func SessionSet(ctx context.Context, key string, value interface{}) (bool, error) {
	token, err := getValidTokenFromContext(ctx)
	if err != nil {
		return false, err
	}

	client, err := data.GetUserServiceClient()
	if err != nil {
		return false, err
	}
	callCtx := fwMd.BuildInnerMDCtx()

	//	connGRPC, callCtx, err := service.NewGrpcClientConn("BossUser")
	//	if err != nil {
	//		return false, err
	//	}
	//
	//	defer connGRPC.Close()
	//
	//	client := userPb.NewUserServiceClient(connGRPC)

	var valueBuf bytes.Buffer
	enc := gob.NewEncoder(&valueBuf)
	err = enc.Encode(value)
	if err != nil {
		return false, err
	}

	reply, err := client.SessionSet(callCtx, &userPb.SessionSetRequest{SessionId: token, Key: key, Value: valueBuf.Bytes()})
	if err != nil {
		return false, err
	}

	return reply.Ok, nil
}

func SessionDelete(ctx context.Context, key string) (bool, error) {
	token, err := getValidTokenFromContext(ctx)
	if err != nil {
		return false, err
	}

	client, err := data.GetUserServiceClient()
	if err != nil {
		return false, err
	}
	callCtx := fwMd.BuildInnerMDCtx()

	//	connGRPC, callCtx, err := service.NewGrpcClientConn("BossUser")
	//	if err != nil {
	//		return false, err
	//	}
	//
	//	defer connGRPC.Close()
	//
	//	client := userPb.NewUserServiceClient(connGRPC)

	reply, err := client.SessionDelete(callCtx, &userPb.SessionDeleteRequest{SessionId: token, Key: key})
	if err != nil {
		return false, err
	}

	return reply.Ok, nil
}

func SessionDestroy(ctx context.Context) (bool, error) {
	token, err := getValidTokenFromContext(ctx)
	if err != nil {
		return false, err
	}

	client, err := data.GetUserServiceClient()
	if err != nil {
		return false, err
	}
	callCtx := fwMd.BuildInnerMDCtx()

	//	connGRPC, callCtx, err := service.NewGrpcClientConn("BossUser")
	//	if err != nil {
	//		return false, err
	//	}
	//
	//	defer connGRPC.Close()
	//
	//	client := userPb.NewUserServiceClient(connGRPC)

	reply, err := client.SessionDestroy(callCtx, &userPb.SessionDestroyRequest{SessionId: token})
	if err != nil {
		return false, err
	}

	return reply.Ok, nil
}
