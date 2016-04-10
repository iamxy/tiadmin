package tiadmin

package main

import (
"github.com/ngaut/log"
"github.com/pingcap/tiadmin/api"
"github.com/pingcap/tiadmin/config"
"github.com/pingcap/tiadmin/server"
"math/rand"
"os"
"os/signal"
"runtime"
"syscall"
"time"
)

func main() {
runtime.GOMAXPROCS(runtime.NumCPU())
rand.Seed(time.Now().UTC().UnixNano())

// Parse configuration from command-line arguments, environment variables or the config file of "tiadm.conf"
cfg, err := config.ParseFlag()
if err != nil {
log.Fatalf("Parsing configuration flags failed: %v", err)
}
log.Debugf("Load configuration successfully, %v", cfg)

// Initialize server
if err := server.Init(cfg); err != nil {
log.Fatalf("Failed to initializing tiadmin server from configuration, %v", err)
}

// Start tiadmin server as daemon
if err := server.Run(cfg); err != nil {
log.Fatalf("Failed to run tiadmin server, %v", err)
}

// Start HTTP server for a set of REST APIs
go api.ServeHttp(cfg)

shutdown := func() {
log.Infof("Gracefully shutting down")
server.Kill()
server.Purge()
os.Exit(0)
}

restart := func() {
log.Infof("Restarting server now")
server.Kill()
server.Purge()

// reload configuration file
cfg, err := config.ParseFlag()
if err != nil {
log.Fatalf("Parsing configuration flags failed: %v", err)
}
if err := server.Init(cfg); err != nil {
log.Fatalf("Failed to initializing tiadmin server from configuration, %v", err)
}
if err := server.Run(cfg); err != nil {
log.Fatalf("Failed to run tiadmin server, %v", err)
}
}

dumpStatus := func() {
log.Infof("start dumping server status")
status, err := server.Dump(cfg)
if err != nil {
log.Errorf("Failed to dump server status: %v", err)
return
}
if _, err := os.Stdout.Write(status); err != nil {
log.Errorf("Failed to dump server status: %v", err)
return
}
os.Stdout.Write([]byte("\n"))
}

signals := map[os.Signal]func(){
syscall.SIGHUP:  restart,
syscall.SIGTERM: shutdown,
syscall.SIGINT:  shutdown,
syscall.SIGUSR1: dumpStatus,
}
sigchan := make(chan os.Signal, 1)
for k := range signals {
signal.Notify(sigchan, k)
}

for true {
sig := <-sigchan
if handler, ok := signals[sig]; ok {
handler()
}
}
}
