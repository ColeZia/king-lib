package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	ke "github.com/go-kratos/kratos/v2/errors"
	mwMetadata "github.com/go-kratos/kratos/v2/middleware/metadata"
	khttp "github.com/go-kratos/kratos/v2/transport/http"

	//"github.com/gogo/protobuf/jsonpb"
	//"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"gl.king.im/king-lib/framework/config"
	fwMd "gl.king.im/king-lib/framework/transport/metadata"
)

type serverConnRegister struct {
	Grpc sync.Map
	Http sync.Map
}

var serverConnReg serverConnRegister

type SvcHttpClient struct {
	*khttp.Client
	ServiceName  string
	DirectlyPath bool
}

type CommonRetWrapper struct {
	//在json.Unmarshal时会影响data的解析，这里先注释
	//ke.Status
	Code     int32             `json:"code"`
	Reason   string            `json:"reason"`
	Message  string            `json:"message"`
	Metadata map[string]string `json:"metadata"`

	Data             interface{} `json:"data"`
	DataProtoMessage proto.Message
	Ts               string `json:"ts"`
	RequestId        string `json:"request_id"`
}

func JsonpbDecoder(ctx context.Context, res *http.Response, out interface{}) error {
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	outWrapper := out.(*CommonRetWrapper)
	err = json.Unmarshal(data, outWrapper)
	if err != nil {
		return err
	}

	outStr, err := json.Marshal(outWrapper.Data)
	if err != nil {
		return err
	}

	//	um := jsonpb.Unmarshaler{
	//		AllowUnknownFields: true,
	//	}
	//	err = um.Unmarshal(strings.NewReader(string(outStr)), outWrapper.DataProtoMessage)

	err = jsonpb.UnmarshalString(string(outStr), outWrapper.DataProtoMessage)
	return err
}

func (c *SvcHttpClient) InvokeWithPathHandle(ctx context.Context, method string, path string, args interface{}, replyData interface{}, opts ...khttp.CallOption) (err error) {
	//先统一加上斜杠
	path = strings.TrimPrefix(path, "/")
	path = "/" + path

	//如果不是直接路径，则加上服务名做路径前缀
	if !c.DirectlyPath {
		path = "/" + c.ServiceName + path
	}

	replywrapper := &CommonRetWrapper{
		Data: replyData,
	}

	err = c.Invoke(ctx, method, path, args, replywrapper, opts...)
	if replywrapper.Code != 200 && replywrapper.Code != 0 {
		statusReason := "InvokeRspError"
		if replywrapper.Reason != "" {
			statusReason = replywrapper.Reason
		}
		kerr := ke.New(int(replywrapper.Code), statusReason, replywrapper.Message)
		err = kerr.WithMetadata(replywrapper.Metadata)
	}

	return
}

func (c *SvcHttpClient) InvokePBWithPathHandle(ctx context.Context, method string, path string, args interface{}, replyData proto.Message, opts ...khttp.CallOption) (err error) {
	//先统一加上斜杠
	path = strings.TrimPrefix(path, "/")
	path = "/" + path

	//如果不是直接路径，则加上服务名做路径前缀
	if !c.DirectlyPath {
		path = "/" + c.ServiceName + path
	}

	replywrapper := &CommonRetWrapper{
		DataProtoMessage: replyData,
	}

	callCtx := fwMd.BuildInnerMDCtx()
	err = c.Invoke(callCtx, method, path, args, replywrapper, opts...)
	if replywrapper.Code != 200 && replywrapper.Code != 0 {
		statusReason := "InvokeRspError"
		if replywrapper.Reason != "" {
			statusReason = replywrapper.Reason
		}
		kerr := ke.New(int(replywrapper.Code), statusReason, replywrapper.Message)
		err = kerr.WithMetadata(replywrapper.Metadata)
	}

	return
}

func NewHttpClientConnInExtNet(serviceName string, cliOpts ...khttp.ClientOption) (*SvcHttpClient, context.Context, error) {
	callCtx := BuildInnerMDCtx()

	if val, ok := serverConnReg.Http.Load(serviceName); ok {
		return val.(*SvcHttpClient), callCtx, nil
	}

	ctx := context.Background()
	//ctx = metadata.AppendToClientContext(ctx, framework.METADATA_KEY_CALL_SCENARIO, framework.MDV_SERVICE_CALL_SCENARIO_INNER)

	serviceConf := config.GetServiceConf()

	var depSvc *config.Service_DependencySvc
	var ok bool
	if depSvc, ok = serviceConf.Service.DependencyHttpServices[serviceName]; !ok {
		log.Println("DependencyServiceNotFound:", serviceName, serviceConf.Service.DependencyHttpServices)
		return nil, nil, errors.New("DependencyServiceNotFound")
	}

	if len(depSvc.Endpoints) < 1 {
		log.Println("DependencyServiceEndpointsEmpty:", serviceName, depSvc)
		return nil, nil, errors.New("DependencyServiceEndpointsEmpty")
	}

	defaultOpts := []khttp.ClientOption{
		khttp.WithEndpoint(depSvc.Endpoints[0]),
		khttp.WithTimeout(time.Second * 10),
		khttp.WithMiddleware(mwMetadata.Client()),
	}

	if depSvc.UseSdnProxy {
		if serviceConf.Service.Sdn == nil || serviceConf.Service.Sdn.ProxyUrl == "" {
			log.Println("SdnProxyUrlEmpty:", serviceName)
			return nil, nil, errors.New("SdnProxyUrlEmpty")
		}

		tr := &http.Transport{
			Proxy: func(_ *http.Request) (*url.URL, error) {
				return url.Parse(serviceConf.Service.Sdn.ProxyUrl)
			},
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			DisableKeepAlives:     false,
			DisableCompression:    false,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ForceAttemptHTTP2:     true,
		}

		defaultOpts = append(defaultOpts, khttp.WithTransport(tr))
	}

	defaultOpts = append(defaultOpts, cliOpts...)
	cliOpts = append(cliOpts)
	kconn, err := khttp.NewClient(
		ctx,
		defaultOpts...,
	)

	conn := &SvcHttpClient{
		Client:      kconn,
		ServiceName: serviceName,
	}

	if depSvc.DirectlyPath {
		conn.DirectlyPath = true
	}

	if err != nil {
		log.Fatal(err)
	}

	serverConnReg.Http.Store(serviceName, conn)

	return conn, callCtx, err
}

func BuildInnerMDCtx() context.Context {
	return fwMd.BuildInnerMDCtx()
}
