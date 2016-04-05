package main

import (
	"github.com/ngaut/log"
	"github.com/pingcap/tiadmin/api"
	"github.com/pingcap/tiadmin/config"
	"github.com/pingcap/tiadmin/server"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse configuration from command-line arguments, environment variables or the config file of "tidbadm.conf"
	cfg, err := config.ParseFlag()
	if err != nil {
		log.Fatalf("Parsing configuration flags failed: %v", err)
	}
	log.SetLevelByString(cfg.LogLevel)

	// Initialize server
	if err := server.Init(cfg); err != nil {
		log.Fatalf("Failed to initializing tidbadm server from configuration, %v", err)
	}

	// Start tidbadm server as daemon
	if err := server.Run(); err != nil {
		log.Fatalf("Starting server unsucessfully, %v", err)
	}

	// Start HTTP server to privide REST APIs
	go api.Serve(cfg)

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
		log.SetLevelByString(cfg.LogLevel)

		if err := server.Init(cfg); err != nil {
			log.Fatalf("Failed to initializing tidbadm server from configuration, %v", err)
		}
		server.Run()
	}

	dumpStatus := func() {
		log.Infof("start dumping server status")
		status, err := server.Dump()
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
