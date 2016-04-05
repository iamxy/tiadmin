package api

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/config"
	"github.com/pingcap/tiadmin/server"
	"net/http"
)

func Serve(cfg *config.Config) {
	if err := setup(cfg); err != nil {
		log.Fatalf("%v", err)
	}
	beego.Run()
}

func setup(cfg *config.Config) error {
	beego.BConfig.AppName = "tidbadm"
	beego.BConfig.RunMode = "dev"
	beego.BConfig.Listen.HTTPPort = 8080

	if err := beegoRouter(); err != nil {
		return err
	}

	beego.InsertFilter("/*", beego.BeforeRouter, func(ctx *context.Context) {
		if !server.IsRunning() {
			ctx.Abort(http.StatusServiceUnavailable, "tidb-admin server is not started up")
		}
	})

	return nil
}
