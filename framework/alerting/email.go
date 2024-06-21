package alerting

import "gl.king.im/king-lib/framework/alerting/types"

var _ NotificationChannel = (*NotiChanWorkEmail)(nil)

type NotiChanWorkEmail struct {
	Key  string
	Name string
}

func (nc *NotiChanWorkEmail) GetKey() string {
	return nc.Key
}

func (nc *NotiChanWorkEmail) GetName() string {
	return nc.Name
}

func (nc *NotiChanWorkEmail) GetType() types.ChannelType {
	return NotiChanKeyEmail
}

func (nc *NotiChanWorkEmail) Notify(message interface{}) bool {
	//todo...
	return true
}

func (nc *NotiChanWorkEmail) NotifyRich(message interface{}) bool {
	//todo...
	return true
}
