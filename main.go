package main

import (
	"flag"
	"io/ioutil"

	"github.com/rpoletaev/test_kpi_drive2/bookshelf"

	"gopkg.in/yaml.v2"
)

func main() {
	init := flag.Bool("init", false, "Create tables and common entities")
	flag.Parse()
	config := &bookshelf.Config{}
	configBytes, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configBytes, config)
	if err != nil {
		panic(err)
	}

	api, erro := bookshelf.NewAPI(*config, *init)
	if erro != nil {
		panic(erro)
	}

	api.Run()
}
