package alerting

import (
	"log"
	"time"

	"gl.king.im/king-lib/framework/alerting/types"
)

const NotiChanKeyWorkWechat types.ChannelType = "WorkWechat"
const NotiChanKeyFeishu types.ChannelType = "Feishu"
const NotiChanKeyEmail types.ChannelType = "Email"

type NotificationChannel interface {
	GetName() (name string)
	GetKey() (key string)
	GetType() (chanType types.ChannelType)
	Notify(message interface{}) bool
	NotifyRich(message interface{}) bool
}

type ChannelGroupRegistry map[string][]string

type Alerting struct {
	ncs         []NotificationChannel
	ncMap       map[string]NotificationChannel
	ncGroup     ChannelGroupRegistry
	msgChan     chan types.Msg
	richMsgChan chan types.ChannelRichMsgMap
}

type AlertFun func(msg interface{}, alerting *Alerting) (succ bool)

//var groupRegistry = ChannelGroupRegistry{}
//
//func RegisterAlertRules(ruleKey string, channels []string) {
//	if _, ok := groupRegistry[ruleKey]; ok {
//		groupRegistry[ruleKey] = append(groupRegistry[ruleKey], channels...)
//	} else {
//		groupRegistry[ruleKey] = channels
//	}
//}

const defaultMsgChanBufferSize = 100

func NewAlerting(ncs []NotificationChannel, ncGroup ChannelGroupRegistry) *Alerting {
	a := &Alerting{
		ncs:         ncs,
		ncMap:       map[string]NotificationChannel{},
		ncGroup:     ncGroup,
		msgChan:     make(chan types.Msg, defaultMsgChanBufferSize),
		richMsgChan: make(chan types.ChannelRichMsgMap, defaultMsgChanBufferSize),
	}

	for _, v := range ncs {
		a.ncMap[v.GetKey()] = v
	}

	a.subMsg()
	a.subRichMsg()

	return a
}

//var DefaultAlerting *Alerting

func (a *Alerting) Alert(msg string) {
	if len(a.ncs) < 1 {
		log.Println("Alerting notification channel empty")
	}

	a.msgChan <- types.Msg{Text: msg}

	//for _, v := range a.ncs {
	//	switch v.GetType() {
	//	case NotiChanKeyWorkWechat, NotiChanKeyEmail:
	//		v.Notify(msg)
	//	}
	//}
}

func (a *Alerting) alertMsg(msg types.Msg) {
	if len(a.ncs) < 1 {
		log.Println("Alerting notification channel empty")
	}

	for _, v := range a.ncs {
		//switch v.GetType() {
		//case NotiChanKeyWorkWechat, NotiChanKeyEmail:
		//	v.Notify(msg.text)
		//}

		v.Notify(msg.Text)
		//防止并发时触发限频
		time.Sleep(400 * time.Millisecond)
	}
}

func (a *Alerting) subMsg() {
	go func() {
		for {
			select {
			case msg := <-a.msgChan:
				a.alertMsg(msg)
			}
		}
	}()
}

func (a *Alerting) subRichMsg() {
	go func() {
		for {
			select {
			case msg := <-a.richMsgChan:
				a.alertRichMsg(msg)
			}
		}
	}()
}

func (a *Alerting) alertRichMsg(msgMap types.ChannelRichMsgMap) {
	if len(a.ncs) < 1 {
		log.Println("Alerting notification channel empty")
	}

	for _, v := range a.ncs {
		//switch v.GetType() {
		//case NotiChanKeyWorkWechat, NotiChanKeyEmail:
		//	v.Notify(msg.text)
		//}

		richMsg, ok := msgMap[v.GetType()]
		if ok {
			v.NotifyRich(richMsg.Content)
		}

		//防止并发时触发限频
		time.Sleep(400 * time.Millisecond)
	}
}

func (a *Alerting) AlertRich(msgMap types.ChannelRichMsgMap) {
	if len(a.ncs) < 1 {
		log.Println("Alerting notification channel empty")
	}

	a.richMsgChan <- msgMap
}

func (a *Alerting) AlertChannel(msg interface{}, channelKey string) bool {
	if nc, ok := a.ncMap[channelKey]; ok {
		return nc.Notify(msg)
	} else {
		log.Println("AlertChannel ncMap has no value for channelKey :" + channelKey)
	}

	return false
}

func (a *Alerting) AlertGroup(msg interface{}, groupKey string) {
	if ncKeys, ok := a.ncGroup[groupKey]; ok {
		if len(ncKeys) < 1 {
			log.Println("AlertGroup ncKeys empty")
		}

		for _, ncKey := range ncKeys {
			if nc, ok := a.ncMap[ncKey]; ok {
				nc.Notify(msg)
			}
		}
	} else {
		log.Println("AlertGroup ncGroup has no value for groupKey :" + groupKey)
	}
}
