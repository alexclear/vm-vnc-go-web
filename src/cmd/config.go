package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Timeout   time.Duration `yaml:"timeout"`
	BindHTTP  string        `yaml:"bind_http"`
	BindHTTPS string        `yaml:"bind_https"`
	Certfile  string
	Keyfile   string
	Debug     bool
}

var cfg = Config{
	Timeout:   time.Minute,
	BindHTTP:  ":80",
	BindHTTPS: ":443",
	Certfile:  "cert.pem",
	Keyfile:   "key.pem",
	Debug:     true,
}

type CommandLineConfig struct {
	config_path *string
}

var commandLineCfg = CommandLineConfig{
	config_path: flag.String("config_path", "config.yaml", "Path to a config file"),
}

func (*CommandLineConfig) Parse() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()

	}

	flag.Parse()
}
