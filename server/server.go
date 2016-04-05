package server

import (
	"errors"
	"fmt"
	etcd "github.com/coreos/etcd/client"
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/agent"
	"github.com/pingcap/tiadmin/config"
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

	ServiceManager  svc.ServiceManager
	AgentReconciler *agent.AgentReconciler
)

func Init(cfg *config.Config) error {
	if IsRunning() {
		return errors.New("Not allowed to initialize a running server")
	}

	// Keep configuration in global scope
	conf = cfg

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
	evtStream := registry.NewEtcdEventStream(kAPI, cfg.EtcdKeyPrefix)

	// Register all Ti-services to a map
	ServiceManager = svc.RegisterServices(reg)

	// Init the manager to monitor and control the local process's state
	procMgr := process.NewProcessManager()

	// Reconciler drives the local process's state towards the desired state
	// stored in the Registry.
	AgentReconciler = agent.NewReconciler(reg, evtStream, procMgr)

	log.Infof("Server initialized successfully")
	return nil
}

func Run() error {
	if IsRunning() {
		return errors.New("Server is already running now")
	}

	// Birth Cry

	stopc = make(chan struct{})
	wg = sync.WaitGroup{}
	components := []func(){
		func() { AgentReconciler.Run(stopc) },
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
}

func Kill() error {
	if !IsRunning() || stopc == nil {
		return errors.New("The server not running, cannot be killed")
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
		return errors.New("Timed out waiting for server to shutdown")
	}

	log.Infof("Tidbadm server stopped")
	switchStateToStopped()
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
