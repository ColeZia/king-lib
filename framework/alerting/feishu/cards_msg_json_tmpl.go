package feishu

import "time"

var tmpl = mapStrIn{
	"config": mapStrIn{
		"wide_screen_mode": true,
	},
	"elements": []mapStrIn{
		{
			"tag":              "column_set",
			"flex_mode":        "none",
			"background_style": "default",
			"columns": []mapStrIn{
				//服务名
				{
					"tag":            "column",
					"width":          "weighted",
					"weight":         1,
					"vertical_align": "top",
					"elements": []mapStrIn{
						{
							"tag": "div",
							"text": mapStrIn{
								"content": "**服务：**",
								"tag":     "lark_md",
							},
						},
					},
				},
				//时间点
				{
					"tag":            "column",
					"width":          "weighted",
					"weight":         1,
					"vertical_align": "top",
					"elements": []mapStrIn{
						{
							"tag": "div",
							"text": mapStrIn{
								"content": "**时间：**" + time.Now().Format(time.RFC3339),
								"tag":     "lark_md",
							},
						},
					},
				},
			},
		},
		//追踪ID
		{
			"tag":     "markdown",
			"content": "**追踪ID：**",
		},
		{
			"tag":              "column_set",
			"flex_mode":        "none",
			"background_style": "default",
			"columns":          []mapStrIn{},
		},
		{
			"tag": "hr",
		},
		//信息文本
		{
			"tag": "div",
			"text": mapStrIn{
				"content": "",
				"tag":     "lark_md",
			},
		},
		{
			"tag": "hr",
		},
		//详细错误信息
		{
			"tag":     "markdown",
			"content": "",
		},
	},
	//标题
	"header": mapStrIn{
		"template": "red",
		"title": mapStrIn{
			"content": "BOSS预警信息",
			"tag":     "plain_text",
		},
	},
}
