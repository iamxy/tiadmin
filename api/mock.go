package api

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/pingcap/tiadmin/schema"
	"math/rand"
	"strconv"
)

var (
	mockProcID = 300000
	mockHosts  = map[string]schema.Host{
		"106001": {
			MachID:   "106001",
			HostName: "mock_node_1",
			HostMeta: schema.HostMeta{Region: "beijing", Datacenter: "DaJiaoTing"},
			PublicIP: "167.124.155.123",
			IsAlive:  true,
			Machine: schema.Machine{
				MachID:     "106001",
				UsageOfCPU: 56,
				TotalMem:   4096,
				UsedMem:    1543,
				TotalSwp:   6010,
				UsedSwp:    4307,
				LoadAvg:    []float32{1.23, 1.67, 1.15},
				UsageOfDisk: []schema.DiskUsage{
					{Mount: "/export", TotalSize: 102400, UsedSize: 5120},
					{Mount: "/", TotalSize: 8192, UsedSize: 578},
				},
			},
		},
		"107001": {
			MachID:   "107001",
			HostName: "mock_node_2",
			HostMeta: schema.HostMeta{Region: "shanghai", Datacenter: "HengBangQiao"},
			PublicIP: "204.224.13.7",
			IsAlive:  true,
			Machine: schema.Machine{
				MachID:     "107001",
				UsageOfCPU: 12,
				TotalMem:   16384,
				UsedMem:    7893,
				TotalSwp:   11034,
				UsedSwp:    2290,
				LoadAvg:    []float32{1.58, 1.91, 1.75},
				UsageOfDisk: []schema.DiskUsage{
					{Mount: "/data", TotalSize: 51200, UsedSize: 38491},
					{Mount: "/", TotalSize: 20480, UsedSize: 7128},
				},
			},
		},
		"108001": {
			MachID:   "108001",
			HostName: "mock_node_3",
			HostMeta: schema.HostMeta{Region: "hongkong", Datacenter: "TongLuoWan"},
			PublicIP: "134.41.99.117",
			IsAlive:  false,
			Machine: schema.Machine{
				MachID:     "108001",
				UsageOfCPU: 0,
				TotalMem:   16384,
				UsedMem:    0,
				TotalSwp:   8319,
				UsedSwp:    0,
				LoadAvg:    []float32{0, 0, 0},
				UsageOfDisk: []schema.DiskUsage{
					{Mount: "/export", TotalSize: 204800, UsedSize: 138910},
					{Mount: "/", TotalSize: 10240, UsedSize: 5549},
				},
			},
		},
		"109001": {
			MachID:   "109001",
			HostName: "mock_node_4",
			HostMeta: schema.HostMeta{Region: "beijing", Datacenter: "TuCheng"},
			PublicIP: "112.94.5.199",
			IsAlive:  true,
			Machine: schema.Machine{
				MachID:     "109001",
				UsageOfCPU: 28,
				TotalMem:   8192,
				UsedMem:    3588,
				TotalSwp:   10241,
				UsedSwp:    6841,
				LoadAvg:    []float32{1.71, 1.54, 1.36},
				UsageOfDisk: []schema.DiskUsage{
					{Mount: "/export", TotalSize: 51200, UsedSize: 1245},
					{Mount: "/", TotalSize: 10960, UsedSize: 7894},
				},
			},
		},
	}
	mockProcs = map[string]schema.Process{
		"356010": {
			ProcID:       "356010",
			SvcName:      "TiDB",
			MachID:       "109001",
			DesiredState: "started",
			CurrentState: "started",
			IsAlive:      false,
			Endpoints:    []string{},
			Executor:     []string{"/bin/sh", "-c"},
			Command:      "tidb-server",
			Args:         []string{"-L", "info", "-path", "/tmp/tidb", "-P", "4000"},
			Environments: []schema.Environment{},
			PublicIP:     "112.94.5.199",
			HostName:     "mock_node_4",
			HostMeta:     schema.HostMeta{Region: "beijing", Datacenter: "TuCheng"},
			Port:         4000,
			Protocol:     "http",
		},
		"356015": {
			ProcID:       "356015",
			SvcName:      "TiDB",
			MachID:       "107001",
			DesiredState: "started",
			CurrentState: "started",
			IsAlive:      true,
			Endpoints:    []string{"http://204.224.13.7:4000"},
			Executor:     []string{"/bin/sh", "-c"},
			Command:      "tidb-server",
			Args:         []string{"-L", "info", "-path", "/tmp/tidb", "-P", "4000"},
			Environments: []schema.Environment{},
			PublicIP:     "204.224.13.7",
			HostName:     "mock_node_2",
			HostMeta:     schema.HostMeta{Region: "shanghai", Datacenter: "HengBangQiao"},
			Port:         4000,
			Protocol:     "http",
		},
		"356020": {
			ProcID:       "356020",
			SvcName:      "TiKV",
			MachID:       "108001",
			DesiredState: "started",
			CurrentState: "started",
			IsAlive:      true,
			Endpoints:    []string{},
			Executor:     []string{"/bin/sh", "-c"},
			Command:      "tikv-server",
			Args:         []string{"-L", "info", "-path", "/tmp/tikv", "-P", "5000"},
			Environments: []schema.Environment{},
			PublicIP:     "134.41.99.117",
			HostName:     "mock_node_3",
			HostMeta:     schema.HostMeta{Region: "hongkong", Datacenter: "TongLuoWan"},
			Port:         5000,
			Protocol:     "http",
		},
	}
	mockServices = map[string]schema.Service{
		"TiDB": {
			SvcName:      "TiDB",
			Version:      "1.0.0",
			Executor:     []string{"/bin/sh", "-c"},
			Command:      "tidb-server",
			Args:         []string{"-L", "info", "-path", "/tmp/tidb", "-P", "4000"},
			Environments: []schema.Environment{},
			Port:         4000,
			Protocol:     "http",
			Dependencies: []string{"TiKV", "PD", "TSO"},
			Endpoints:    []string{"http://204.224.13.7:4000"},
		},
		"TiKV": {
			SvcName:      "TiKV",
			Version:      "1.0.0",
			Executor:     []string{"/bin/sh", "-c"},
			Command:      "tikv-server",
			Args:         []string{"-L", "info", "-path", "/tmp/tikv", "-P", "5000"},
			Environments: []schema.Environment{},
			Port:         5000,
			Protocol:     "http",
			Dependencies: []string{"TSO", "PD"},
			Endpoints:    []string{},
		},
		"TSO": {
			SvcName:      "TSO",
			Version:      "1.0.0",
			Executor:     []string{"/bin/sh", "-c"},
			Command:      "tso-server",
			Args:         []string{"-L", "info"},
			Environments: []schema.Environment{},
			Port:         4100,
			Protocol:     "http",
			Dependencies: []string{},
			Endpoints:    []string{},
		},
		"PD": {
			SvcName:      "PD",
			Version:      "1.0.0",
			Executor:     []string{"/bin/sh", "-c"},
			Command:      "pd-server",
			Args:         []string{"-L", "info"},
			Environments: []schema.Environment{},
			Port:         4200,
			Protocol:     "http",
			Dependencies: []string{},
			Endpoints:    []string{},
		},
	}
)

func mockRouter() error {
	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/version", &MockVersionController{}, "get:VersionInfo"),
		beego.NSRouter("/hosts", &MockHostController{}, "get:FindAllHosts"),
		beego.NSRouter("/hosts/:machID", &MockHostController{}, "get:FindHost"),
		beego.NSRouter("/hosts/:machID/meta", &MockHostController{}, "put:SetHostMetaInfo"),
		beego.NSRouter("/services", &MockServiceController{}, "get:AllServices"),
		beego.NSRouter("/services/:svcName", &MockServiceController{}, "get:Service"),
		beego.NSRouter("/processes", &MockProcessController{}, "get:FindAllProcesses"),
		beego.NSRouter("/processes", &MockProcessController{}, "post:StartNewProcess"),
		beego.NSRouter("/processes/findByHost", &MockProcessController{}, "get:FindByHost"),
		beego.NSRouter("/processes/findByService", &MockProcessController{}, "get:FindByService"),
		beego.NSRouter("/processes/:procID", &MockProcessController{}, "get:FindProcess"),
		beego.NSRouter("/processes/:procID", &MockProcessController{}, "delete:DestroyProcess"),
		beego.NSRouter("/processes/:procID/start", &MockProcessController{}, "get:StartProcess"),
		beego.NSRouter("/processes/:procID/stop", &MockProcessController{}, "get:StopProcess"),
		beego.NSRouter("/monitor/real/tidb_perf", &MockMonitorController{}, "get:TiDBPerformanceMetrics"),
		beego.NSRouter("/monitor/real/tikv_storage", &MockMonitorController{}, "get:TiKVStorageMetrics"),
	)
	beego.AddNamespace(ns)
	return nil
}

type MockVersionController struct {
	beego.Controller
}

type MockHostController struct {
	beego.Controller
}

type MockServiceController struct {
	beego.Controller
}

type MockProcessController struct {
	beego.Controller
}

type MockMonitorController struct {
	beego.Controller
}

func (c *MockVersionController) VersionInfo() {
	c.Data["json"] = schema.Version{
		Version:      "1.0.0",
		BuildUTCTime: "2016-01-19 08:12:47",
	}
}

func (c *MockHostController) FindAllHosts() {
	hosts := []schema.Host{}
	for _, val := range mockHosts {
		hosts = append(hosts, val)
	}
	c.Data["json"] = hosts
	c.ServeJSON()
}

func (c *MockHostController) FindHost() {
	machID := c.Ctx.Input.Param(":machID")
	if len(machID) == 0 {
		c.Abort("400")
	}
	host, ok := mockHosts[machID]
	if !ok {
		c.Abort("404")
	}
	c.Data["json"] = host
	c.ServeJSON()
}

func (c *MockHostController) SetHostMetaInfo() {
	machID := c.Ctx.Input.Param(":machID")
	if len(machID) == 0 {
		c.Abort("400")
	}
	host, ok := mockHosts[machID]
	if !ok {
		c.Abort("404")
	}
	var meta schema.HostMeta
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &meta)
	if err != nil || len(meta.Region) == 0 || len(meta.Datacenter) == 0 {
		c.Abort("400")
	}
	host.HostMeta = meta
	mockHosts[machID] = host
	c.Data["json"] = host
	c.ServeJSON()
}

func (c *MockServiceController) AllServices() {
	services := []schema.Service{}
	for _, val := range mockServices {
		services = append(services, val)
	}
	c.Data["json"] = services
	c.ServeJSON()
}

func (c *MockServiceController) Service() {
	svcName := c.Ctx.Input.Param(":svcName")
	if len(svcName) == 0 {
		c.Abort("400")
	}
	svc, ok := mockServices[svcName]
	if !ok {
		c.Abort("404")
	}
	c.Data["json"] = svc
	c.ServeJSON()
}

func (c *MockProcessController) FindAllProcesses() {
	procs := []schema.Process{}
	for _, val := range mockProcs {
		procs = append(procs, val)
	}
	c.Data["json"] = procs
	c.ServeJSON()
}

func (c *MockProcessController) StartNewProcess() {
	var proc schema.Process
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &proc)
	if err != nil || len(proc.SvcName) == 0 || len(proc.MachID) == 0 {
		c.Abort("400")
	}
	if svc, ok := mockServices[proc.SvcName]; ok {
		if len(proc.Executor) == 0 {
			proc.Executor = svc.Executor
		}
		if len(proc.Command) > 0 {
			proc.Command = svc.Command
		}
		if len(proc.Args) == 0 {
			proc.Args = svc.Args
		}
		if len(proc.Environments) == 0 {
			proc.Environments = svc.Environments
		}
		if proc.Port == 0 {
			proc.Port = svc.Port
		}
		if len(proc.Protocol) == 0 {
			proc.Protocol = svc.Protocol
		}
	} else {
		c.Abort("400")
	}
	if mach, ok := mockHosts[proc.MachID]; ok {
		proc.PublicIP = mach.PublicIP
		proc.HostName = mach.HostName
		proc.HostMeta = mach.HostMeta
	} else {
		c.Abort("400")
	}
	proc.ProcID = strconv.Itoa(mockProcID)
	mockProcID++
	if len(proc.DesiredState) == 0 {
		proc.DesiredState = "started"
	}
	proc.CurrentState = "stopped"
	proc.IsAlive = false
	mockProcs[proc.ProcID] = proc
	c.Data["json"] = proc
	c.ServeJSON()
}

func (c *MockProcessController) FindByHost() {
	machID := c.GetString("machID")
	if len(machID) == 0 {
		c.Abort("400")
	}
	if _, ok := mockHosts[machID]; ok {
		res := []schema.Process{}
		for _, proc := range mockProcs {
			if proc.MachID == machID {
				res = append(res, proc)
			}
		}
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		c.Abort("404")
	}
}

func (c *MockProcessController) FindByService() {
	svcName := c.GetString("svcName")
	if len(svcName) == 0 {
		c.Abort("400")
	}
	if _, ok := mockServices[svcName]; ok {
		res := []schema.Process{}
		for _, proc := range mockProcs {
			if proc.SvcName == svcName {
				res = append(res, proc)
			}
		}
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		c.Abort("404")
	}
}

func (c *MockProcessController) FindProcess() {
	procID := c.Ctx.Input.Param(":procID")
	if len(procID) == 0 {
		c.Abort("400")
	}
	if proc, ok := mockProcs[procID]; ok {
		c.Data["json"] = proc
		c.ServeJSON()
	} else {
		c.Abort("404")
	}
}

func (c *MockProcessController) DestroyProcess() {
	procID := c.Ctx.Input.Param(":procID")
	if len(procID) == 0 {
		c.Abort("400")
	}
	if proc, ok := mockProcs[procID]; ok {
		delete(mockProcs, procID)
		c.Data["json"] = proc
		c.ServeJSON()
	} else {
		c.Abort("404")
	}
}

func (c *MockProcessController) StartProcess() {
	procID := c.Ctx.Input.Param(":procID")
	if len(procID) == 0 {
		c.Abort("400")
	}
	if proc, ok := mockProcs[procID]; ok {
		proc.DesiredState = "started"
		mockProcs[procID] = proc
		c.Data["json"] = proc
		c.ServeJSON()
	} else {
		c.Abort("404")
	}
}

func (c *MockProcessController) StopProcess() {
	procID := c.Ctx.Input.Param(":procID")
	if len(procID) == 0 {
		c.Abort("400")
	}
	if proc, ok := mockProcs[procID]; ok {
		proc.DesiredState = "stopped"
		mockProcs[procID] = proc
		c.Data["json"] = proc
		c.ServeJSON()
	} else {
		c.Abort("404")
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func (c *MockMonitorController) TiDBPerformanceMetrics() {
	c.Data["json"] = schema.PerfMetrics{
		Tps:  int32(randInt(20, 100)),
		Qps:  int32(randInt(65, 300)),
		Iops: int32(randInt(80, 550)),
	}
	c.ServeJSON()
}

func (c *MockMonitorController) TiKVStorageMetrics() {
	c.Data["json"] = schema.StorageMetrics{
		Usage:    int64(randInt(535, 565)),
		Capacity: 8192,
	}
	c.ServeJSON()
}
