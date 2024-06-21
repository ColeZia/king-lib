package alerting

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"

	"gl.king.im/king-lib/framework/alerting/types"
)

var _ NotificationChannel = (*NotiChanFeishu)(nil)

type NotiChanFeishu struct {
	Key            string
	Name           string
	Webhook        string
	ContentType    string
	ContentBuilder func(msg interface{}) string
	withTrace      bool
	Debug          bool
}

func (nc *NotiChanFeishu) GetKey() string {
	return nc.Key
}

func (nc *NotiChanFeishu) GetName() string {
	return nc.Name
}

func (nc *NotiChanFeishu) GetType() types.ChannelType {
	return NotiChanKeyFeishu
}

func Notifiy() {

}
func (nc *NotiChanFeishu) Notify(message interface{}) bool {
	if message == nil {
		return true
	}

	var content string
	if nc.ContentBuilder != nil {
		content = nc.ContentBuilder(message)
	} else {

		if nc.withTrace {
			//content += fmt.Sprintf("`%s`\n", alertMsg)
			stack := debug.Stack()
			content += fmt.Sprintf("%s\n\nStack::%s", message, stack)
		} else {
			content += fmt.Sprintf("%s\n", message)
		}

		//暂时保持和企微一致
		if len(content) > 4000 {
			content = content[0:4000]
		}
	}

	contentJson := map[string]string{
		"text": content,
	}

	contentJsonBytes, err := json.Marshal(contentJson)

	if err != nil {
		fmt.Println("NotiChanFeishu post body marshal err::", err)
	}

	if nc.Debug {
		fmt.Println("NotiChanFeishu contentJsonStr:", string(contentJsonBytes))
	}

	postBody := map[string]interface{}{
		"msg_type": "text",
		"content":  string(contentJsonBytes),
	}

	return nc.notify(postBody)
}

func (nc *NotiChanFeishu) NotifyRich(message interface{}) bool {
	if message == nil {
		return true
	}

	var content string
	if nc.ContentBuilder != nil {
		content = nc.ContentBuilder(message)
	} else {

		if nc.withTrace {
			//content += fmt.Sprintf("`%s`\n", alertMsg)
			stack := debug.Stack()
			content += fmt.Sprintf("%s\n\nStack::%s", message, stack)
		} else {
			content += fmt.Sprintf("%s\n", message)
		}

		//暂时保持和企微一致
		if len(content) > 4000 {
			content = content[0:4000]
		}
	}

	contentJson := map[string]string{
		"text": content,
	}

	_ = contentJson

	contentJsonBytes, err := json.Marshal(message)

	if err != nil {
		fmt.Println("NotiChanFeishu post body marshal err::", err)
	}

	if nc.Debug {
		fmt.Println("NotiChanFeishu contentJsonStr:", string(contentJsonBytes))
	}

	postBody := map[string]interface{}{
		"msg_type": "interactive",
		"card":     string(contentJsonBytes),
	}

	return nc.notify(postBody)
}

func (nc *NotiChanFeishu) notify(postBody map[string]interface{}) bool {

	postBodyJson, err := json.Marshal(postBody)

	if err != nil {
		fmt.Println("NotiChanFeishu post body marshal err::", err)
	}

	go func(postJsonBytes []byte) {
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				fmt.Println("NotiChanFeishu post recover::", recoverErr)
			}
		}()

		if nc.Debug {
			fmt.Println("NotiChanFeishu post:", nc.Webhook, string(postJsonBytes))
		}

		resp, err := http.Post(nc.Webhook, "application/json", strings.NewReader(string(postJsonBytes)))

		if err != nil {
			fmt.Println("alert post err", err)
		}

		if resp != nil && nc.Debug {
			if resp.StatusCode != http.StatusOK {
				fmt.Println("alert post resp.StatusCode", resp.StatusCode)
			}

			defer resp.Body.Close()
			//转写body
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("alert post resp err", err)
			}
			fmt.Println("alert post resp body::", string(bodyBytes))
		}

	}(postBodyJson)
	//_ = resp
	//if err != nil {
	//	panic(err)
	//}

	//if resp.StatusCode != http.StatusOK {
	//	panic("warning http status " + strconv.Itoa(resp.StatusCode))
	//}

	return true
}
