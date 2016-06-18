package main

import (
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

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

	http.Handle("/metrics", prometheus.Handler())

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
