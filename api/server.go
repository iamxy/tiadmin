package api

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	//"github.com/astaxie/beego/plugins/cors"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/config"
	"github.com/pingcap/tiadmin/frontend"
	"github.com/pingcap/tiadmin/schema"
	"github.com/pingcap/tiadmin/server"
)

func bad_request(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	modelError := &schema.ModelError{
		ErrCode: http.StatusBadRequest,
		Reason:  "Bad request parameters",
	}
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

	//beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
	//	AllowOrigins:     []string{"http://localhost:9000"},
	//	AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
	//	ExposeHeaders:    []string{"Content-Length"},
	//	AllowCredentials: true,
	//}))

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
				ctx.Abort(http.StatusServiceUnavailable, "tiadmin server is not available")
			}
		})
	}

	// router for static resources
	beego.Handler("/*", http.FileServer(
		&assetfs.AssetFS{
			Asset:     frontend.Asset,
			AssetDir:  frontend.AssetDir,
			AssetInfo: frontend.AssetInfo,
		},
	))
	beego.Run()
}
