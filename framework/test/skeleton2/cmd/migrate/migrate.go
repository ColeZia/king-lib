/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package migrate

import (
	"context"
	"fmt"
	"log"
	"os"

	"gl.king.im/king-lib/framework/test/skeleton2/internal/conf"
	"gl.king.im/king-lib/framework/test/skeleton2/internal/data/ent/demo"

	"gl.king.im/king-lib/framework/service/cmd/migrate"
)

var Bc = &conf.Bootstrap{}

func MigrateRunFun(migrateMode string) {

	client := newAuthClient(Bc.Data)
	defer client.Close()
	ctx := context.Background()
	var err error
	//表结构同步
	opts := migrate.DefaultOpts()
	switch migrateMode {
	case "print":
		log.Println("迁移开始-print")
		err = client.Debug().Schema.WriteTo(ctx, os.Stdout, opts...)
	case "sync":
		log.Println("迁移开始-sync")

		err = client.Debug().Schema.Create(ctx, opts...)
	default:
		log.Fatalf("不支持的迁移模式:" + migrateMode)
	}

	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	log.Println("迁移完成--" + migrateMode)
}

func newAuthClient(c *conf.Data) *demo.Client {
	dsn := fmt.Sprintf("%s?charset=utf8mb4&parseTime=True&loc=Local", c.Database.Source)
	opts := []demo.Option{
		//user.Debug(),
	}
	cli, err := demo.Open("mysql", dsn, opts...)
	if err != nil {
		panic(err)
	}
	return cli
}
