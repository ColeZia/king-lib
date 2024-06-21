package data

import (
	"fmt"

	"gl.king.im/king-lib/framework/test/skeleton/internal/conf"
	"gl.king.im/king-lib/framework/test/skeleton/internal/data/ent/demo"

	"git.e.coding.king.cloud/dev/data_platform/gorm"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewBossDb,
	NewBossEntCli,
	NewSkeletonRepo,
)

type BossDb *gorm.DB

// Data .
type Data struct {
	bossDb     BossDb
	bossEntCli *demo.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, bossDb BossDb, bossEntCli *demo.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{bossDb: bossDb, bossEntCli: bossEntCli}, cleanup, nil
}

func NewBossDb(c *conf.Data) BossDb {

	return BossDb(nil)
}

func NewBossEntCli(c *conf.Data) *demo.Client {
	dsn := fmt.Sprintf("%s?charset=utf8mb4&parseTime=True&loc=Local", c.Database.Source)
	cli, err := demo.Open("mysql", dsn, demo.Debug())
	if err != nil {
		panic(err)
	}
	return cli
}
