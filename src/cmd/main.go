package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"

	"runtime/pprof"

	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

func sighandler() {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)
	if cfg.Debug {
		pprof.StopCPUProfile()
	}
	os.Exit(0)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func main() {
	commandLineCfg.Parse()

	content, err := ioutil.ReadFile(*commandLineCfg.config_path)
	if err != nil {
		log.Errorf("Error reading configuration file: %v", err)
	}

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		log.Fatalf("Error parsing configuration file: %v", err)
	}

	spew.Printf("%#+v\n%#+v\n", commandLineCfg, cfg)

	if cfg.Debug {
		f, err := os.Create("server-go.prof")
		if err != nil {
			log.Fatal("Error enabling CPU profiler: ", err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	go sighandler()

	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", httpHandler)

	go func() {
		log.Debugf("Serving...")
		// err = http.ListenAndServeTLS(cfg.BindHTTPS, cfg.Certfile, cfg.Keyfile, nil)
		err = http.ListenAndServe(cfg.BindHTTPS, nil)
		if err != nil {
			log.Fatal("ListenAndServeTLS: ", err)
		}
	}()

	err = http.ListenAndServe(cfg.BindHTTP, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
