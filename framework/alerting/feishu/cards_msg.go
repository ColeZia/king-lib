package feishu

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"gl.king.im/king-lib/framework/alerting/types"
	"gl.king.im/king-lib/framework/service/appinfo"
)

const MsgTypeCard = "card"

type mapStrIn map[string]interface{}

type BossAlertCard interface {
	BuildMsg(mapStrIn) mapStrIn
}

type bossAlertCardSvcVals struct {
	SvcName  string
	Content  string
	TraceID  string
	Details  string
	ID       string
	Elements []Element
}

type BossAlertCardVals struct {
	Title    string
	Content  string
	Details  string
	TraceID  string
	Elements []Element
}

type Text struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type ElementColumn struct {
	Tag           string    `json:"tag"`
	Width         string    `json:"width"`
	Weight        int       `json:"weight"`
	VerticalAlign string    `json:"vertical_align"`
	Elements      []Element `json:"elements"`
}

type Element struct {
	Tag             string          `json:"tag"`
	FlexMode        string          `json:"flex_mode"`
	BackgroundStyle string          `json:"background_style"`
	Columns         []ElementColumn `json:"columns"`
	Content         string          `json:"content"`
	Text            Text            `json:"text"`
}

type CardHeader struct {
	Template string            `json:"template"`
	Title    map[string]string `json:"title"`
}

type CardMsg struct {
	Config   map[string]bool `json:"config"`
	Elements []Element       `json:"elements"`
	Header   CardHeader      `json:"header"`
}

func buildBossAlertCardSvcElements(vals bossAlertCardSvcVals) []Element {
	var elements []Element
	elements = append(elements, Element{
		Tag:             "column_set",
		FlexMode:        "none",
		BackgroundStyle: "default",
		Columns: []ElementColumn{
			{
				Tag:           "column",
				Width:         "weighted",
				Weight:        1,
				VerticalAlign: "top",
				Elements: []Element{
					{
						Tag: "div",
						Text: Text{
							Tag:     "lark_md",
							Content: "**服务: **" + vals.SvcName,
						},
					},
				},
			},
			{
				Tag:           "column",
				Width:         "weighted",
				Weight:        1,
				VerticalAlign: "top",
				Elements: []Element{
					{
						Tag: "div",
						Text: Text{
							Tag:     "lark_md",
							Content: "**时间: **" + time.Now().Format("2006-01-02 15:04:05"),
						},
					},
				},
			},
		},
	})

	//追踪ID
	elements = append(elements, Element{
		//Tag:     "markdown",
		//Content: "**追踪ID: **" + vals.TraceID,

		Tag: "div",
		Text: Text{
			Tag:     "lark_md",
			Content: "**追踪ID: **" + vals.TraceID,
		},
	})

	//服务ID
	elements = append(elements, Element{
		//Tag:     "markdown",
		//Content: "**实例ID: **" + vals.ID,

		Tag: "div",
		Text: Text{
			Tag:     "lark_md",
			Content: "**实例ID: **" + vals.ID,
		},
	})

	//elements = append(elements, Element{
	//	Tag:             "column_set",
	//	FlexMode:        "none",
	//	BackgroundStyle: "default",
	//})

	elements = append(elements, Element{
		Tag: "hr",
	})

	if len(vals.Elements) > 0 {
		elements = append(elements, vals.Elements...)
		elements = append(elements, Element{
			Tag: "hr",
		})
	}

	//信息文本
	elements = append(elements, Element{
		Tag: "div",
		Text: Text{
			Content: vals.Content,
			Tag:     "lark_md",
		},
	})

	elements = append(elements, Element{
		Tag: "hr",
	})

	elements = append(elements, Element{
		Tag:     "markdown",
		Content: vals.Details,
	})
	return elements
}

func buildBaseSvcVals(ctx context.Context, vals BossAlertCardVals) bossAlertCardSvcVals {
	traceIdValuer := tracing.TraceID()
	traceId := traceIdValuer(ctx).(string)

	if vals.TraceID != "" {
		traceId = vals.TraceID
	}

	svcVals := bossAlertCardSvcVals{
		SvcName:  appinfo.AppInfoIns.Name,
		ID:       appinfo.AppInfoIns.Id,
		TraceID:  traceId,
		Content:  vals.Content,
		Details:  vals.Details,
		Elements: vals.Elements,
	}

	return svcVals
}

func BossAlertCardRichMsg(ctx context.Context, alertLev types.AlertLevel, vals BossAlertCardVals) (richMsg types.RichMsg) {
	switch alertLev {
	case types.AlertLevelError:
		richMsg = BossAlertCardSvcError(ctx, vals)
	case types.AlertLevelInfo:
		richMsg = BossAlertCardSvcInfo(ctx, vals)
	case types.AlertLevelWarn:
		richMsg = BossAlertCardSvcWarn(ctx, vals)
	default:
		richMsg = BossAlertCardSvcInfo(ctx, vals)
	}

	return
}

func BossAlertCardSvcError(ctx context.Context, vals BossAlertCardVals) (richMsg types.RichMsg) {
	title := "BOSS预警信息"
	if vals.Title != "" {
		title = vals.Title
	}
	cardMsg := CardMsg{
		Header: CardHeader{
			Template: "red",
			Title: map[string]string{
				"content": "【错误】" + title,
				"tag":     "plain_text",
			},
		},
		Config: map[string]bool{
			"wide_screen_mode": true,
		},
		Elements: buildBossAlertCardSvcElements(buildBaseSvcVals(ctx, vals)),
	}

	richMsg = types.RichMsg{Type: MsgTypeCard, Content: cardMsg}

	return
}

func BossAlertCardSvcInfo(ctx context.Context, vals BossAlertCardVals) (richMsg types.RichMsg) {
	title := "BOSS预警信息"
	if vals.Title != "" {
		title = vals.Title
	}

	cardMsg := CardMsg{
		Header: CardHeader{
			Template: "blue",
			Title: map[string]string{
				"content": "【信息】" + title,
				"tag":     "plain_text",
			},
		},
		Config: map[string]bool{
			"wide_screen_mode": true,
		},
		Elements: buildBossAlertCardSvcElements(buildBaseSvcVals(ctx, vals)),
	}

	richMsg = types.RichMsg{Type: MsgTypeCard, Content: cardMsg}

	return
}

func BossAlertCardSvcWarn(ctx context.Context, vals BossAlertCardVals) (richMsg types.RichMsg) {
	title := "BOSS预警信息"
	if vals.Title != "" {
		title = vals.Title
	}

	cardMsg := CardMsg{
		Header: CardHeader{
			Template: "orange",
			Title: map[string]string{
				"content": "【警告】" + title,
				"tag":     "plain_text",
			},
		},
		Config: map[string]bool{
			"wide_screen_mode": true,
		},
		Elements: buildBossAlertCardSvcElements(buildBaseSvcVals(ctx, vals)),
	}

	richMsg = types.RichMsg{Type: MsgTypeCard, Content: cardMsg}

	return
}

func BossAlertCardSvcSuccess(ctx context.Context, vals BossAlertCardVals) (richMsg types.RichMsg) {
	title := "BOSS预警信息"
	if vals.Title != "" {
		title = vals.Title
	}

	cardMsg := CardMsg{
		Header: CardHeader{
			Template: "green",
			Title: map[string]string{
				"content": "【成功】" + title,
				"tag":     "plain_text",
			},
		},
		Config: map[string]bool{
			"wide_screen_mode": true,
		},
		Elements: buildBossAlertCardSvcElements(buildBaseSvcVals(ctx, vals)),
	}

	richMsg = types.RichMsg{Type: MsgTypeCard, Content: cardMsg}

	return
}

func BossAlertCardSvcExpired(ctx context.Context, vals BossAlertCardVals) (richMsg types.RichMsg) {

	title := "BOSS预警信息"
	if vals.Title != "" {
		title = vals.Title
	}

	cardMsg := CardMsg{
		Header: CardHeader{
			Template: "grey",
			Title: map[string]string{
				"content": "【过期】" + title,
				"tag":     "plain_text",
			},
		},
		Config: map[string]bool{
			"wide_screen_mode": true,
		},
		Elements: buildBossAlertCardSvcElements(buildBaseSvcVals(ctx, vals)),
	}

	richMsg = types.RichMsg{Type: MsgTypeCard, Content: cardMsg}

	return
}
