package api

import (
	"github.com/astaxie/beego"
)

func beegoRouter() error {
	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/version", &VersionController{}, "get:VersionInfo"),
		beego.NSRouter("/hosts", &HostController{}, "get:FindAllHosts"),
		beego.NSRouter("/hosts/:machID", &HostController{}, "get:FindHost"),
		beego.NSRouter("/hosts/:machID/meta", &HostController{}, "put:SetHostMetaInfo"),
		beego.NSRouter("/services", &ServiceController{}, "get:AllServices"),
		beego.NSRouter("/services/:svcName", &ServiceController{}, "get:Service"),
		beego.NSRouter("/processes", &ProcessController{}, "get:FindAllProcesses"),
		beego.NSRouter("/processes", &ProcessController{}, "post:StartNewProcess"),
		beego.NSRouter("/processes/findByHost", &ProcessController{}, "get:FindByHost"),
		beego.NSRouter("/processes/findByService", &ProcessController{}, "get:FindByService"),
		beego.NSRouter("/processes/:procID", &ProcessController{}, "get:FindProcess"),
		beego.NSRouter("/processes/:procID", &ProcessController{}, "delete:DestroyProcess"),
		beego.NSRouter("/processes/:procID/start", &ProcessController{}, "get:StartProcess"),
		beego.NSRouter("/processes/:procID/stop", &ProcessController{}, "get:StopProcess"),
		beego.NSRouter("/monitor/real/tidb_perf", &MonitorController{}, "get:TiDBPerformanceMetrics"),
		beego.NSRouter("/monitor/real/tikv_storage", &MonitorController{}, "get:TiKVStorageMetrics"),
	)
	beego.AddNamespace(ns)
	return nil
}
