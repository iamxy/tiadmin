package server

import (
	"errors"
	"fmt"
	etcd "github.com/coreos/etcd/client"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/agent"
	"github.com/pingcap/tiadmin/config"
	"github.com/pingcap/tiadmin/machine"
	"github.com/pingcap/tiadmin/process"
	"github.com/pingcap/tiadmin/registry"
	svc "github.com/pingcap/tiadmin/service"
	"sync"
	"time"
)

const (
	shutdownTimeout = time.Minute
)

var (
	stopc   chan struct{}  // used to terminate all other goroutines
	wg      sync.WaitGroup // used to co-ordinate shutdown
	running bool           = false
	conf    *config.Config

	SvcMgr   svc.ServiceManager
	AgentRec *agent.AgentReconciler
)

func Init(cfg *config.Config) error {
	if IsRunning() {
		return errors.New("Not allowed to initialize a running server")
	}

	// Keep configuration in global scope
	conf = cfg

	agentTTL, err := time.ParseDuration(cfg.AgentTTL)
	if err != nil {
		return err
	}

	// Init registry driver of etcd, and event stream
	etcdRequestTimeout := time.Duration(cfg.EtcdRequestTimeout) * time.Millisecond
	etcdCfg := etcd.Config{
		Endpoints: cfg.EtcdServers,
		Transport: etcd.DefaultTransport,
	}
	etcdClient, err := etcd.New(etcdCfg)
	if err != nil {
		return err
	}
	kAPI := etcd.NewKeysAPI(etcdClient)
	reg := registry.NewEtcdRegistry(kAPI, cfg.EtcdKeyPrefix, etcdRequestTimeout)
	eStream := registry.NewEtcdEventStream(kAPI, cfg.EtcdKeyPrefix)

	// Register all Ti-services to a map
	SvcMgr = svc.RegisterServices(reg)

	// Init process manager for local processes
	procMgr := process.NewProcessManager()
	// Init machine
	mach := machine.NewMachine()
	// Create agent
	ag := agent.New(reg, procMgr, mach, agentTTL)

	// Reconciler drives the local process's state towards the desired state
	// stored in the Registry.
	AgentRec = agent.NewReconciler(reg, eStream, ag)

	log.Infof("Server initialized successfully")
	return nil
}

func Run() (err error) {
	if IsRunning() {
		err = errors.New("Server is already running, cannot call to run repeatly")
		return
	}
	if conf.IsMock {
		log.Warnf("Server is started in mock mode, skip to Run()")
		return
	}

	// Birth Cry

	stopc = make(chan struct{})
	wg = sync.WaitGroup{}
	components := []func(){
		func() { AgentRec.Run(stopc) },
	}

	for _, f := range components {
		f := f
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	log.Infof("Server started successfully")
	switchStateToRunning()
	return
}

func Kill() (err error) {
	if !IsRunning() || stopc == nil {
		err = errors.New("The server is not running, cannot be killed")
		return
	}

	close(stopc)
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(shutdownTimeout):
		err = errors.New("Timed out waiting for server to shutdown")
		return
	}

	log.Infof("Tiadmin server stopped")
	switchStateToStopped()
	return
}

func Purge() {
}

func Dump() (dumpinfo []byte, err error) {
	err = nil
	dumpinfo = []byte(fmt.Sprintf("%v", conf))
	log.Infof("Finished dumping server status")
	return
}

func IsRunning() bool {
	return running
}

func switchStateToStopped() {
	running = false
}

func switchStateToRunning() {
	running = true
}
