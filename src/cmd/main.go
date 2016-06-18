package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
)

func main() {
	commandLineCfg.Parse()

	content, err := ioutil.ReadFile(*commandLineCfg.config_path)
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		log.Fatalf("Error parsing configuration file: %v", err)
	}

	go func() {
		log.Debugf("Serving...")
		err = http.ListenAndServeTLS(cfg.BindHTTPS, cfg.Certfile, cfg.Keyfile, nil)
		if err != nil {
			incFatalErrorsCount(fatalErrorListenAndServeTLS)
			dumpFatalErrors()
			log.Fatal("ListenAndServeTLS: ", err)
		}
		wg.Done()
	}()

	spew.Printf("%#+v\n%#+v\n", commandLineCfg, cfg)
}
