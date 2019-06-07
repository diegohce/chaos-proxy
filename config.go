package main

import (
	"encoding/json"
	"io/ioutil"
)

type hostConfig struct {
	Host     string `json:"host"`
	RealPath string `json:"real_path"`
}

type proxyConfig struct {
	DefaultHost hostConfig `json:"default_host"`
	Paths       map[string] hostConfig `json:"paths"`
}

var (
	chaosConfig proxyConfig
)

func loadConfig() error {

	if b, err := ioutil.ReadFile("./chaos-proxy.json"); err != nil {
		return err
	} else {
		log.Debug().Println("config content", string(b))
		return json.Unmarshal(b, &chaosConfig)
	}
}
