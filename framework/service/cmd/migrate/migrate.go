/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package migrate

import (
	"context"
	"fmt"

	"gl.king.im/king-lib/framework/service/appinfo"
	"gl.king.im/king-lib/framework/service/cmd/config"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/casbin/ent-adapter/ent/migrate"
	"github.com/spf13/cobra"
)

type MigrateRunFun func(mode string)
type SeedRunFun func()

var (
	flagconf      string
	drop          = false
	migrateMode   = "print"
	migrateRunFun MigrateRunFun
	seedRunFun    SeedRunFun
	MigrateBc     interface{}
)

// serviceCmd represents the service command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "数据迁移相关命令",
	Long:  `数据迁移相关命令，包括数据迁移执行等功能`,
	Run: func(cmd *cobra.Command, args []string) {
		config.ScanCnf(appinfo.AppInfoIns.Name, flagconf, MigrateBc)
		migrateRunFun(migrateMode)
	},
}

// serviceCmd represents the service command
var migrateSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "微服务初始化命令",
	Long:  `微服务初始化项目代码`,
	Run: func(cmd *cobra.Command, args []string) {
		config.ScanCnf(appinfo.AppInfoIns.Name, flagconf, MigrateBc)
		seedRunFun()
	},
}

func InitCmd(rootCmd *cobra.Command, mf MigrateRunFun, sf SeedRunFun, mbc interface{}) {
	MigrateBc = mbc
	migrateRunFun = mf
	seedRunFun = sf

	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	migrateCmd.PersistentFlags().StringVarP(&flagconf, "conf", "c", "", "配置文件")
	migrateCmd.PersistentFlags().StringVarP(&migrateMode, "mode", "m", "print", "迁移模式")
	migrateCmd.PersistentFlags().BoolVar(&drop, "drop", false, "是否执行删除操作")

}

func DefaultOpts() []schema.MigrateOption {
	opts := []schema.MigrateOption{}

	hooksOpt := schema.WithHooks(func(next schema.Creator) schema.Creator {
		return schema.CreateFunc(func(ctx context.Context, tables ...*schema.Table) error {

			finalTables := []*schema.Table{}
			for _, v := range tables {

				//if v.Name == "xxx" {
				//	continue
				//}

				finalTables = append(finalTables, v)

				fmt.Println("table:", v.Name)
			}
			// Run custom code here.
			return next.Create(ctx, finalTables...)
		})
	})
	opts = append(opts, hooksOpt)

	//关闭外键
	opts = append(opts, migrate.WithForeignKeys(false))
	if drop {
		opts = append(opts, migrate.WithDropIndex(true))
		opts = append(opts, migrate.WithDropColumn(true))
	}

	return opts
}
