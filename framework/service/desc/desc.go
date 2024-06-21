package desc

import "google.golang.org/protobuf/types/descriptorpb"

type ServiceDescMethodMap struct {
	ServDesc  *descriptorpb.ServiceDescriptorProto
	MethodMap map[string]*descriptorpb.MethodDescriptorProto
}

var (
	ServerRegCnfMap map[string]ServiceDescMethodMap
)

func init() {
	ServerRegCnfMap = map[string]ServiceDescMethodMap{}
}
