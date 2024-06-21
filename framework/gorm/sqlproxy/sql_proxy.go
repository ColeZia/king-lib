package sqlproxy

import (
	"git.e.coding.king.cloud/dev/data_platform/gorm"
	"gl.king.im/king-lib/framework/auth/user"
	"gl.king.im/king-lib/protobuf/api/common"
)

const dpPrefix = "data_platform:"
const feishuPrefix = "feishu:"

type sqlProxy struct{}

var _ gorm.Plugin = (*sqlProxy)(nil)

func (p *sqlProxy) Initialize(db *gorm.DB) error {
	db.Callback().Create().Before("*").Register("gorm:sql_proxy", proxyCallback)
	db.Callback().Query().Before("*").Register("gorm:sql_proxy", proxyCallback)
	db.Callback().Update().Before("*").Register("gorm:sql_proxy", proxyCallback)
	db.Callback().Delete().Before("*").Register("gorm:sql_proxy", proxyCallback)
	db.Callback().Row().Before("*").Register("gorm:sql_proxy", proxyCallback)
	db.Callback().Raw().Before("*").Register("gorm:sql_proxy", proxyCallback)
	return nil
}

func (p *sqlProxy) Name() string {
	return "gorm:sql_proxy"
}

func proxyCallback(tx *gorm.DB) {
	ctx := tx.Statement.Context
	token, ok := user.GlobalOceanTokenFromServerContext(ctx)
	if !ok {
		global, exists := user.OpUserGlobalInfoFromServerContext(ctx)
		if !exists {
			return
		}
		switch global.SessionInfo.LoginMethod {
		case common.LoginMethod_Ocean:
			token = dpPrefix + global.SessionInfo.OceanAuthToken
		case common.LoginMethod_Feishu:
			token = feishuPrefix + global.SessionInfo.FeishuAuthToken
		}
	}
	tx.Token(token)
}

func New() gorm.Plugin {
	return &sqlProxy{}
}
