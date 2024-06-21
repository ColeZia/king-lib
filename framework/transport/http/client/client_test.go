package client

import (
	"context"
	"log"
	"reflect"
	"testing"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	fwMd "gl.king.im/king-lib/framework/transport/metadata"
	"gl.king.im/king-lib/protobuf/api/common"
	userPb "gl.king.im/king-lib/protobuf/api/user/service/v1"
)

func TestNewHttpClientConnInExtNet(t *testing.T) {
	type args struct {
		serviceName string
		cliOpts     []khttp.ClientOption
	}
	tests := []struct {
		name    string
		args    args
		want    *SvcHttpClient
		want1   context.Context
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			conn, callCtx, err := NewHttpClientConnInExtNet(common.SERVICE_NAME_BossUser.String())
			if err != nil {
				return
			}
			_ = callCtx

			args := &userPb.AuthCheckRequest{Token: "k-0M3h9MKPj99mvAgPI388Q5lagQuxTiNJfbvxxEqKU=",
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
				log.Println("", "client.AuthCheck err::", err)
				return
			}

			got, got1, err := NewHttpClientConnInExtNet(tt.args.serviceName, tt.args.cliOpts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHttpClientConnInExtNet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHttpClientConnInExtNet() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("NewHttpClientConnInExtNet() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
