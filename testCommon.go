package main

import (
	"apiserver/config"
	"testing"

	"github.com/spf13/pflag"
)

var cfgTest = pflag.StringP("config", "t", "", "apiserver config file path.")

func InitTestingConfig(m *testing.M) {
	pflag.Parse()

	// init config
	if err := config.Init(*cfgTest); err != nil {
		panic(err)
	}
}
