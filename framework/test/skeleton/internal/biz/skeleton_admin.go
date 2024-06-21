package biz

import (
	"context"
	"fmt"

	"gl.king.im/king-lib/framework/auth/user"
	"gl.king.im/king-lib/framework/log"
	ssUser "gl.king.im/king-lib/framework/session/user"
	"gl.king.im/king-lib/framework/transport/http/client"
	"gl.king.im/king-lib/framework/transport/metadata"

	ke "github.com/go-kratos/kratos/v2/errors"
	v1 "gl.king.im/king-lib/framework/test/skeleton2/api/service/v1"
	pb "gl.king.im/king-lib/protobuf/api/skeleton/admin/v1"

	fwMd "gl.king.im/king-lib/framework/transport/metadata"
	"gl.king.im/king-lib/protobuf/api/common"
	userPb "gl.king.im/king-lib/protobuf/api/user/service/v1"
)

//框架骨架示例repository

func (uc *SkeletonUsecase) AdminAuthGet(ctx context.Context, req *pb.GetRequest) (rep *pb.GetReply, err error) {

	entity, err := ssUser.GetUserIdentity(ctx)
	if err != nil {
		return
	}

	_ = entity

	globalInfo, ok := user.OpUserGlobalInfoFromServerContext(ctx)
	if ok {
		fmt.Println("globalInfo::", globalInfo.OceanAuthToken)
	}

	log.ChangeStackLogModule(ctx, "jjiiioooo")
	log.ChangeStackLogModule(ctx, "xxxx")
	uc.log.Infom(ctx, "gggggg")
	uc.log.Warnm(ctx, "mmmmmm")
	return
	cli, err := NewTest2Client()

	_, err = cli.Get(ctx, &v1.GetRequest{})
	return

	conn, _, err := client.NewHttpClientConnInExtNet(common.SERVICE_NAME_BossUser.String())
	if err != nil {
		return
	}

	args := &userPb.AuthCheckRequest{Token: "_Ti9SJSU5FmazeQ9H-HmAuY9iwv5753_5yzOapapzqg=",
		Resource:             "/ffdf",
		AuthenticationIgnore: false,
		AuthorizationIgnore:  false}

	callCtx1 := fwMd.BuildInnerMDCtx()

	//type commonWrapper struct {
	//	Code      int
	//	Data      interface{}
	//	Ts        string
	//	Message   string
	//	RequestId string
	//}

	type HttpBossOperationUser struct {
		Id       uint64               `protobuf:"varint,100,opt,name=id,proto3" json:"id,string,omitempty"`
		Username string               `protobuf:"bytes,200,opt,name=username,proto3" json:"username,omitempty"`
		Nickname string               `protobuf:"bytes,300,opt,name=nickname,proto3" json:"nickname,omitempty"`
		Name     string               `protobuf:"bytes,400,opt,name=name,proto3" json:"name,omitempty"` //姓名
		Mobile   string               `protobuf:"bytes,500,opt,name=mobile,proto3" json:"mobile,omitempty"`
		Gender   string               `protobuf:"varint,600,opt,name=gender,proto3,enum=common.BossOperationUser_Gender" json:"gender,omitempty"`
		Email    string               `protobuf:"bytes,700,opt,name=email,proto3" json:"email,omitempty"`   //邮箱
		Avatar   string               `protobuf:"bytes,800,opt,name=avatar,proto3" json:"avatar,omitempty"` //头像
		Employee *common.BossEmployee `protobuf:"bytes,900,opt,name=employee,proto3" json:"employee,omitempty"`
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
		fmt.Println("", "client.AuthCheck err::", err)
		return
	}

	return

	//panic("AdminAuthGet...")

	//err = errors.ErrorStarPicSubCurrencyNotAllowed("菜单不存在")
	//err = ke.InternalServer("TEST", "TEST...")
	return
	//cli, err := NewPayClient()

	callCtx := metadata.BuildInnerMDCtx()
	_ = callCtx

	//_, err = cli.GetPayOrderListByOrders(ctx, &v1.GetPayOrderListByOrdersRequest{})

	//uc.log.Info("Info aaa")
	//uc.log.Infoc(ctx, "Infoc bbb")
	//uc.log.Infos(ctx, "Infos-key1", 1111, "Infos-key2", 2222, 3333, 33331111)
	//uc.log.Infow("Infow-key1", 111, "Infow-key2")
	panic("ggg")
	klh := uc.log.GetKratosLogHelper()

	klh.Info("original kratos helper:::")

	uc.log.Debuga(ctx, "a-debug", 1111, "a-key2", 2222)
	uc.log.Infoc(context.Background(), "c-info", 1111, "a-key2", 2222)
	uc.log.Warn(ctx, "warn", 1111, "key2", 2222)
	uc.log.Errora(ctx, "a-error", 1111, "a-key2", 2222)

	uc.log.Infom(ctx, "m-key1", 1111, "m-key2")
	return

	lh := uc.log.WithContext(ctx)
	lh.Info("(uc *SkeletonUsecase) Get...", "val111", "key2", "val22")
	lh.Infow("key111", "val111", "key2", "val22")

	loh, ok := log.LogHelperFromServerContext(ctx)
	if ok {
		loh.Debug("LogHelperFromServerContext log")
	}

	err = ke.InternalServer("TEST", "TEST-MSG")
	//err = ke.BadRequest("", "")

	return
}
