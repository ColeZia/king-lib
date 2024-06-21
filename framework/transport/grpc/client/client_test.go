package client

import (
	"reflect"
	"testing"

	"google.golang.org/grpc"
)

func TestNewGrpcClientConn(t *testing.T) {
	type args struct {
		serviceName string
		opts        []newGrpcClientConnOption
	}
	tests := []struct {
		name    string
		args    args
		want    *grpc.ClientConn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGrpcClientConn(tt.args.serviceName, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGrpcClientConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGrpcClientConn() = %v, want %v", got, tt.want)
			}
		})
	}
}
