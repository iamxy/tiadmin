package api

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/config"
	"github.com/pingcap/tiadmin/schema"
	"github.com/pingcap/tiadmin/server"
	"net/http"
)

func bad_request(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	modelError := schema.ModelError{ErrCode: http.StatusBadRequest, Reason: "Bad request parameters"}
	json, err := json.Marshal(modelError)
	if err != nil {
		log.Errorf("Failed sending HTTP response body: %v", err)
	}
	_, err = rw.Write(json)
	if err != nil {
		log.Errorf("Failed sending HTTP response body: %v", err)
	}
}

func ServeHttp(cfg *config.Config) {
	beego.BConfig.AppName = "tiadmin"
	beego.BConfig.RunMode = "dev"
	beego.BConfig.Listen.HTTPPort = 8080
	beego.BConfig.CopyRequestBody = true

	beego.ErrorHandler("400", bad_request)

	if cfg.IsMock {
		if err := mockRouter(); err != nil {
			log.Fatalf("parsing beego router error, %v", err)
		}
	} else {
		if err := beegoRouter(); err != nil {
			log.Fatalf("parsing beego router error, %v", err)
		}
		beego.InsertFilter("/*", beego.BeforeRouter, func(ctx *context.Context) {
			if !server.IsRunning() {
				ctx.Abort(http.StatusServiceUnavailable, "tiadmin server not ready now")
			}
		})
	}

	beego.Run()
}
