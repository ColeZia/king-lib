package alerting

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"gl.king.im/king-lib/framework/alerting/types"
)

var _ NotificationChannel = (*NotiChanWorkWechat)(nil)

type NotiChanWorkWechat struct {
	Key            string
	Name           string
	Webhook        string
	ContentType    string
	ContentBuilder func(msg interface{}) string
	withTrace      bool
	Debug          bool
}

func (nc *NotiChanWorkWechat) GetKey() string {
	return nc.Key
}

func (nc *NotiChanWorkWechat) GetName() string {
	return nc.Name
}

func (nc *NotiChanWorkWechat) GetType() types.ChannelType {
	return NotiChanKeyWorkWechat
}

func (nc *NotiChanWorkWechat) Notify(message interface{}) bool {

	if message != nil {
		var content string
		if nc.ContentBuilder != nil {
			content = nc.ContentBuilder(message)
		} else {
			content = ""

			if nc.withTrace {
				//content += fmt.Sprintf("`%s`\n", alertMsg)
				stack := debug.Stack()
				content += fmt.Sprintf("%s\n\nStack::%s", message, stack)
			} else {
				content += fmt.Sprintf("%s\n", message)
			}

			//内容过长会报错，这里做截取，报错提示最长4096，这里保险一点也同时也为了避免信息过长取4000
			//{"errcode":40058,"errmsg":"markdown.content exceed max length 4096. Invalid input invalid Request Parameter, hint: [1647424338252370220001484], from ip: 113.116.157.231, more info at https://open.work.weixin.qq.com/devtool/query?e=40058"}
			if len(content) > 4000 {
				content = content[0:4000]
			}
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

		go func(postJsonBytes []byte) {
			defer func() {
				if recoverErr := recover(); recoverErr != nil {
					fmt.Println("WorkWechatAlert post recover::", recoverErr)
				}
			}()

			resp, err := http.Post(nc.Webhook, "application/json", strings.NewReader(string(postJsonBytes)))

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

	return true
}

func (nc *NotiChanWorkWechat) NotifyRich(message interface{}) bool {
	//暂时复用基本notify
	return nc.Notify(message)
}
