package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"gl.king.im/king-lib/framework/alerting"
	"gl.king.im/king-lib/framework/alerting/feishu"
	"gl.king.im/king-lib/framework/alerting/types"
	"gl.king.im/king-lib/framework/config"
	"gl.king.im/king-lib/framework/internal/di"
	"gl.king.im/king-lib/framework/service/appinfo"
)

func Alert(alertMsg interface{}, opts ...interface{}) {
	AlertPlain(context.Background(), alertMsg)
}

func AlertPlain(ctx context.Context, alertMsg interface{}) {
	if di.Container != nil && di.Container.DefaultAlerting != nil {
		di.Container.DefaultAlerting.Alert(fmt.Sprintf("服务:%s\n时间:%s\n\n%+v", appinfo.AppInfoIns.Name, time.Now().Format("2006-01-02 15:04:05"), alertMsg))
	} else {
		log.Println("无默认预警器实例")
	}
}

func AlertRich(ctx context.Context, alertMsg interface{}, alertLev types.AlertLevel) {
	if di.Container != nil && di.Container.DefaultAlerting != nil {

		feishuVals := feishu.BossAlertCardVals{
			Title:   "预警消息",
			Content: fmt.Sprint(alertMsg),
			Details: "",
			//自定义字段
			//Elements: []feishu.Element{
			//	{
			//		Tag: "div",
			//		Text: feishu.Text{
			//			Tag:     "lark_md",
			//			Content: "**Operation: **" + "abc/efg/hij",
			//		},
			//	},
			//	{
			//		Tag: "div",
			//		Text: feishu.Text{
			//			Tag:     "lark_md",
			//			Content: "**User: **" + "test",
			//		},
			//	},
			//},
		}

		feishuRichMsg := feishu.BossAlertCardRichMsg(ctx, alertLev, feishuVals)

		richMsgMap := types.ChannelRichMsgMap{
			//飞书
			alerting.NotiChanKeyFeishu: feishuRichMsg,
			//企微
			//alerting.NotiChanKeyWorkWechat: types.RichMsg{Content: "企微通知消息内容..."},
		}

		di.Container.DefaultAlerting.AlertRich(richMsgMap)
	} else {
		log.Println("无默认预警器实例")
	}
}

func AlertWithTrace(alertMsg interface{}) {
	WorkWechatAlert(alertMsg, true)
}

func WorkWechatAlert(alertMsg interface{}, opts ...interface{}) {

	if alertMsg != nil {
		var withTrace bool
		//打印追踪信息
		if len(opts) > 0 {
			withTrace = opts[0].(bool)
		}

		if withTrace {
			//printAllTrace := false
			//if len(opts) > 1 {
			//	printAllTrace = opts[1].(bool)
			//}
			//fmt.Println("alertMsg::", alertMsg)
			//stack := make([]byte, 1<<16)
			//_ = printAllTrace
			//runtime.Stack(stack, printAllTrace)
			//fmt.Println("runtime stack::", string(stack))
		}

		serviceConf := config.GetServiceConf()
		if serviceConf.Service.Alert == nil || serviceConf.Service.Alert.WorkWechat == nil || serviceConf.Service.Alert.WorkWechat.Hook == "" {
			fmt.Println("Service.Alert.WorkWechat.Hook配置为空！")
			return
		}

		content := "alert msg::"

		if withTrace {
			//content += fmt.Sprintf("`%s`\n", alertMsg)
			stack := debug.Stack()
			content += fmt.Sprintf("`%s`\n\nStack::%s", alertMsg, stack)
		} else {
			content += fmt.Sprintf("`%s`\n", alertMsg)
		}

		//内容过长会报错，这里做截取，报错提示最长4096，这里保险一点也同时也为了避免信息过长取4000
		//{"errcode":40058,"errmsg":"markdown.content exceed max length 4096. Invalid input invalid Request Parameter, hint: [1647424338252370220001484], from ip: 113.116.157.231, more info at https://open.work.weixin.qq.com/devtool/query?e=40058"}
		if len(content) > 4000 {
			content = content[0:4000]
		}

		postBody := map[string]interface{}{
			"msgtype": "markdown",
			"markdown": map[string]string{
				"content": content,
			},
		}

		postBodyJson, err := json.Marshal(postBody)

		if err != nil {
			fmt.Println("WorkWechatAlert post body marshal err::", err)
		}

		if serviceConf.Service.Env == "local" {
			return //本地关闭
		}

		go func(postJsonBytes []byte) {
			defer func() {
				if recoverErr := recover(); recoverErr != nil {
					fmt.Println("WorkWechatAlert post recover::", recoverErr)
				}
			}()

			resp, err := http.Post(serviceConf.Service.Alert.WorkWechat.Hook, "application/json", strings.NewReader(string(postJsonBytes)))

			if err != nil {
				fmt.Println("alert post err", err)
			}

			if resp != nil {
				if resp.StatusCode != http.StatusOK {
					fmt.Println("alert post resp.StatusCode", resp.StatusCode)
				}

				//defer resp.Body.Close()
				////转写body
				//bodyBytes, err := ioutil.ReadAll(resp.Body)
				//if err != nil {
				//	fmt.Println("alert post resp err", err)
				//}
				//fmt.Println("alert post resp body::", string(bodyBytes))
			}

		}(postBodyJson)
		//_ = resp
		//if err != nil {
		//	panic(err)
		//}

		//if resp.StatusCode != http.StatusOK {
		//	panic("warning http status " + strconv.Itoa(resp.StatusCode))
		//}
	}

}
