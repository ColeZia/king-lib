package alert

import (
	"errors"
	"fmt"

	httpstatus "github.com/go-kratos/kratos/v2/transport/http/status"
	"gl.king.im/king-lib/framework/alerting"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 500
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

type AlertError struct {
	Code     int32             `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Reason   string            `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	Message  string            `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	Metadata map[string]string `protobuf:"bytes,4,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Alerting *alerting.Alerting
	AlertFun alerting.AlertFun
}

func (x *AlertError) String() string {
	return x.Message
}

func (x *AlertError) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *AlertError) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *AlertError) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *AlertError) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (e *AlertError) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v", e.Code, e.Reason, e.Message, e.Metadata)
}

// GRPCStatus returns the Status represented by se.
func (e *AlertError) GRPCStatus() *status.Status {
	s, _ := status.New(httpstatus.ToGRPCCode(int(e.Code)), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   e.Reason,
			Metadata: e.Metadata,
		})
	return s
}

// Is matches each error in the chain with the target value.
func (e *AlertError) Is(err error) bool {
	if se := new(AlertError); errors.As(err, &se) {
		return se.Reason == e.Reason
	}
	return false
}
