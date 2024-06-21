package alerting

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gl.king.im/king-lib/framework/alerting/feishu"
	"gl.king.im/king-lib/framework/alerting/types"
)

func TestAlerting_Alert(t *testing.T) {
	type fields struct {
		ncs     []NotificationChannel
		ncMap   map[string]NotificationChannel
		ncGroup ChannelGroupRegistry
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "case1:并发测试",
			fields: fields{
				ncs: []NotificationChannel{&NotiChanFeishu{
					Debug:   true,
					Key:     "payAlert",
					Name:    "",
					Webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/768da0f1-5d78-4828-b429-2b737e68bbef",
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAlerting(tt.fields.ncs, nil)

			//富文本json嵌套序列化问题测试
			if true {
				str := `aaaaa{\n      \"id\": \"pi_3N14ZiKL1kL8OjRP1o2C3sJ8\",\n      \"object\": \"payment_intent\",\n      \"amount\": 2000,\n      \"amount_capturable\": 0,\n      \"amount_details\": {\n        \"tip\": {\n        }\n      },\n      \"amount_received\": 0,\n      \"application\": null,\n      \"application_fee_amount\": null,\n      \"automatic_payment_methods\": null,\n      \"canceled_at\": null,\n      \"cancellation_reason\": null,\n      \"capture_method\": \"automatic\",\n      \"charges\": {\n        \"object\": \"list\",\n        \"data\": [\n          {\n            \"id\": \"ch_3N14ZiKL1kL8OjRP1mDan71o\",\n            \"object\": \"charge\",\n            \"amount\": 2000,\n            \"amount_captured\": 0,\n            \"amount_refunded\": 0,\n            \"application\": null,\n            \"application_fee\": null,\n            \"application_fee_amount\": null,\n            \"balance_transaction\": null,\n            \"billing_details\": {\n              \"address\": {\n                \"city\": null,\n                \"country\": null,\n                \"line1\": null,\n                \"line2\": null,\n                \"postal_code\": null,\n                \"state\": null\n              },\n              \"email\": null,\n              \"name\": null,\n              \"phone\": null\n            },\n            \"calculated_statement_descriptor\": \"Stripe\",\n            \"captured\": false,\n            \"created\": 1682501027,\n            \"currency\": \"usd\",\n            \"customer\": null,\n            \"description\": \"(created by Stripe CLI)\",\n            \"destination\": null,\n            \"dispute\": null,\n            \"disputed\": false,\n            \"failure_balance_transaction\": null,\n            \"failure_code\": \"card_declined\",\n            \"failure_message\": \"Your card was declined.\",\n            \"fraud_details\": {\n            },\n            \"invoice\": null,\n            \"livemode\": false,\n            \"metadata\": {\n            },\n            \"on_behalf_of\": null,\n            \"order\": null,\n            \"outcome\": {\n              \"network_status\": \"declined_by_network\",\n              \"reason\": \"generic_decline\",\n              \"risk_level\": \"normal\",\n              \"risk_score\": 37,\n              \"seller_message\": \"The bank did not return any further details with this decline.\",\n              \"type\": \"issuer_declined\"\n            },\n            \"paid\": false,\n            \"payment_intent\": \"pi_3N14ZiKL1kL8OjRP1o2C3sJ8\",\n            \"payment_method\": \"pm_1N14ZiKL1kL8OjRPIBO4BrwY\",\n            \"payment_method_details\": {\n              \"card\": {\n                \"brand\": \"visa\",\n                \"checks\": {\n                  \"address_line1_check\": null,\n                  \"address_postal_code_check\": null,\n                  \"cvc_check\": null\n                },\n                \"country\": \"US\",\n                \"exp_month\": 4,\n                \"exp_year\": 2024,\n                \"fingerprint\": \"gphiiwOP1JeI02mu\",\n                \"funding\": \"credit\",\n                \"installments\": null,\n                \"last4\": \"0002\",\n                \"mandate\": null,\n                \"network\": \"visa\",\n                \"network_token\": {\n                  \"used\": false\n                },\n                \"three_d_secure\": null,\n                \"wallet\": null\n              },\n              \"type\": \"card\"\n            },\n            \"receipt_email\": null,\n            \"receipt_number\": null,\n            \"receipt_url\": null,\n            \"refunded\": false,\n            \"refunds\": {\n              \"object\": \"list\",\n              \"data\": [\n\n              ],\n              \"has_more\": false,\n              \"total_count\": 0,\n              \"url\": \"/v1/charges/ch_3N14ZiKL1kL8OjRP1mDan71o/refunds\"\n            },\n            \"review\": null,\n            \"shipping\": null,\n            \"source\": null,\n            \"source_transfer\": null,\n            \"statement_descriptor\": null,\n            \"statement_descriptor_suffix\": null,\n            \"status\": \"failed\",\n            \"transfer_data\": null,\n            \"transfer_group\": null\n          }\n        ],\n        \"has_more\": false,\n        \"total_count\": 1,\n        \"url\": \"/v1/charges?payment_intent=pi_3N14ZiKL1kL8OjRP1o2C3sJ8\"\n      },\n      \"client_secret\": \"pi_3N14ZiKL1kL8OjRP1o2C3sJ8_secret_a9QDk4nafBoGC1MGr98NA292N\",\n      \"confirmation_method\": \"automatic\",\n      \"created\": 1682501026,\n      \"currency\": \"usd\",\n      \"customer\": null,\n      \"description\": \"(created by Stripe CLI)\",\n      \"invoice\": null,\n      \"last_payment_error\": {\n        \"charge\": \"ch_3N14ZiKL1kL8OjRP1mDan71o\",\n        \"code\": \"card_declined\",\n        \"decline_code\": \"generic_decline\",\n        \"doc_url\": \"https://stripe.com/docs/error-codes/card-declined\",\n        \"message\": \"Your card was declined.\",\n        \"payment_method\": {\n          \"id\": \"pm_1N14ZiKL1kL8OjRPIBO4BrwY\",\n          \"object\": \"payment_method\",\n          \"billing_details\": {\n            \"address\": {\n              \"city\": null,\n              \"country\": null,\n              \"line1\": null,\n              \"line2\": null,\n              \"postal_code\": null,\n              \"state\": null\n            },\n            \"email\": null,\n            \"name\": null,\n            \"phone\": null\n          },\n          \"card\": {\n            \"brand\": \"visa\",\n            \"checks\": {\n              \"address_line1_check\": null,\n              \"address_postal_code_check\": null,\n              \"cvc_check\": null\n            },\n            \"country\": \"US\",\n            \"exp_month\": 4,\n            \"exp_year\": 2024,\n            \"fingerprint\": \"gphiiwOP1JeI02mu\",\n            \"funding\": \"credit\",\n            \"generated_from\": null,\n            \"last4\": \"0002\",\n            \"networks\": {\n              \"available\": [\n                \"visa\"\n              ],\n              \"preferred\": null\n            },\n            \"three_d_secure_usage\": {\n              \"supported\": true\n            },\n            \"wallet\": null\n          },\n          \"created\": 1682501026,\n          \"customer\": null,\n          \"livemode\": false,\n          \"metadata\": {\n          },\n          \"type\": \"card\"\n        },\n        \"type\": \"card_error\"\n      },\n      \"latest_charge\": \"ch_3N14ZiKL1kL8OjRP1mDan71o\",\n      \"livemode\": false,\n      \"metadata\": {\n      },\n      \"next_action\": null,\n      \"on_behalf_of\": null,\n      \"payment_method\": null,\n      \"payment_method_options\": {\n        \"card\": {\n          \"installments\": null,\n          \"mandate_options\": null,\n          \"network\": null,\n          \"request_three_d_secure\": \"automatic\"\n        }\n      },\n      \"payment_method_types\": [\n        \"card\"\n      ],\n      \"processing\": null,\n      \"receipt_email\": null,\n      \"review\": null,\n      \"setup_future_usage\": null,\n      \"shipping\": null,\n      \"source\": null,\n      \"statement_descriptor\": null,\n      \"statement_descriptor_suffix\": null,\n      \"status\": \"requires_payment_method\",\n      \"transfer_data\": null,\n      \"transfer_group\": null\n    }`

				richMsgMap := types.ChannelRichMsgMap{
					NotiChanKeyFeishu: feishu.BossAlertCardSvcInfo(context.Background(), feishu.BossAlertCardVals{
						Title:   "测试标题222",
						Content: str,
						Details: "飞书通知消息详情...",
						Elements: []feishu.Element{
							{
								Tag: "div",
								Text: feishu.Text{
									Tag:     "lark_md",
									Content: "**Operation: **" + "abc/efg/hij",
								},
							},
							{
								Tag: "div",
								Text: feishu.Text{
									Tag:     "lark_md",
									Content: "**User: **" + "test",
								},
							},
						}},
					),
					NotiChanKeyWorkWechat: types.RichMsg{Content: "企微通知消息内容..."},
				}

				a.AlertRich(richMsgMap)
				time.Sleep(3 * time.Second)
			}

			//富文本消息测试块
			if false {

				a.AlertRich(types.ChannelRichMsgMap{NotiChanKeyFeishu: feishu.BossAlertCardSvcError(context.Background(), feishu.BossAlertCardVals{Title: "测试标题111", Content: "测试内容44", Details: "stack list..."})})
				richMsgMap := types.ChannelRichMsgMap{
					NotiChanKeyFeishu: feishu.BossAlertCardSvcInfo(context.Background(), feishu.BossAlertCardVals{
						Title:   "测试标题222",
						Content: "飞书通知消息内容",
						Details: "飞书通知消息详情...",
						Elements: []feishu.Element{
							{
								Tag: "div",
								Text: feishu.Text{
									Tag:     "lark_md",
									Content: "**Operation: **" + "abc/efg/hij",
								},
							},
							{
								Tag: "div",
								Text: feishu.Text{
									Tag:     "lark_md",
									Content: "**User: **" + "test",
								},
							},
						}},
					),
					NotiChanKeyWorkWechat: types.RichMsg{Content: "企微通知消息内容..."},
				}

				a.AlertRich(richMsgMap)
				a.AlertRich(types.ChannelRichMsgMap{NotiChanKeyFeishu: feishu.BossAlertCardSvcWarn(context.Background(), feishu.BossAlertCardVals{Title: "测试标题333", Content: "测试内容66", Details: "stack list..."})})
				a.AlertRich(types.ChannelRichMsgMap{NotiChanKeyFeishu: feishu.BossAlertCardSvcExpired(context.Background(), feishu.BossAlertCardVals{Title: "测试标题444", Content: "测试内容77", Details: "stack list..."})})
				a.AlertRich(types.ChannelRichMsgMap{NotiChanKeyFeishu: feishu.BossAlertCardSvcSuccess(context.Background(), feishu.BossAlertCardVals{Title: "测试标题555", Content: "测试内容88", Details: "stack list..."})})

				time.Sleep(5 * time.Second)
			}

			//并发测试块
			if false {

				testFlag := "case1:并发测试-5:"
				for i := 0; i < 5; i++ {
					go func(i int) {
						a.Alert(fmt.Sprintf(testFlag+"通知A:%d", i))
					}(i)
				}

				for i := 0; i < 5; i++ {
					go func(i int) {
						a.Alert(fmt.Sprintf(testFlag+"通知B:%d", i))
					}(i)
				}

				a.Alert(testFlag + "并发通知")
				time.Sleep(time.Second)
				a.Alert(testFlag + "并发后通知")

				time.Sleep(10 * time.Second)
			}

		})
	}
}
