package config

//type ServiceConfig struct {
//	state         protoimpl.MessageState
//	sizeCache     protoimpl.SizeCache
//	unknownFields protoimpl.UnknownFields
//
//	Server struct {
//		state         protoimpl.MessageState
//		sizeCache     protoimpl.SizeCache
//		unknownFields protoimpl.UnknownFields
//		Http          *struct {
//			state         protoimpl.MessageState
//			sizeCache     protoimpl.SizeCache
//			unknownFields protoimpl.UnknownFields
//
//			Network string               `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
//			Addr    string               `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
//			Timeout *durationpb.Duration `protobuf:"bytes,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
//		}
//		Grpc *struct {
//			state         protoimpl.MessageState
//			sizeCache     protoimpl.SizeCache
//			unknownFields protoimpl.UnknownFields
//
//			Network string               `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
//			Addr    string               `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
//			Timeout *durationpb.Duration `protobuf:"bytes,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
//		}
//	}
//
//	//Server  *Server  `protobuf:"bytes,1,opt,name=server,proto3" json:"server,omitempty"`
//	Data struct {
//		Database struct {
//			Driver string
//			Source string
//		}
//		Redis struct {
//			Network      string
//			Addr         string
//			ReadTimeout  *durationpb.Duration
//			WriteTimeout *durationpb.Duration
//		}
//	}
//
//	Service struct {
//		Registrys struct {
//			Consul struct {
//				Addr string
//			}
//			Etcd struct {
//				Addr string
//			}
//		}
//
//		Traces struct {
//			Jaeger struct {
//				Endpoint string
//			}
//		}
//
//		Alert struct {
//			WorkWechat struct {
//				Hook string
//			} `json:"work_wechat"`
//		}
//
//		Env string
//	}
//}

var AppConfIns = &Bootstrap{} //ServiceConfig
var ConfSource string

//func SetServiceConf(ac *conf.Bootstrap) {
//	AppConfIns = ac
//}

func GetServiceConf() *Bootstrap {
	//暂时取消重新加载
	//func GetServiceConf(opts ...interface{}) conf.Bootstrap {

	//	confSrc := ConfSource
	//	if len(opts) > 0 {
	//		confSrc = opts[0].(string)
	//	}
	//
	//	if confSrc == "" {
	//		panic("conf source is empty!")
	//	}
	//
	//	LoadServiceConf(confSrc)
	return AppConfIns
}

//应在框架启动流程的起点执行
func LoadServiceConf(s string) {
	ConfSource = s
	c := GetInstanceBySource(ConfSource)
	// Unmarshal the config to struct
	if err := c.Scan(AppConfIns); err != nil {
		panic(err)
	}
}
