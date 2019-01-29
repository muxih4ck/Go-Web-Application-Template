package main

import (
	"apiserver/config"
	"apiserver/model"
	"apiserver/router"
	"apiserver/router/middleware"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

var (
	cfgTest  = pflag.StringP("config_test", "t", "", "apiserver config file path.")
	g        *gin.Engine
	token    string
	username string
	password string
	uid      uint64
)


// This function is used to do setup before executing the test functions
func TestMain(m *testing.M) {

	pflag.Parse()

	// init config
	if err := config.Init(*cfgTest); err != nil {
		panic(err)
	}
	fmt.Println(viper.GetString("db.password"))
	fmt.Println(viper.GetString("db.addr"))
	// init db
	model.DB.Init()
	defer model.DB.Close()
	//Set Gin to Test Mode
	gin.SetMode(viper.GetString("runmode"))

	g = gin.New()
	router.Load(
		// Cores.
		g,

		// Middlwares.
		middleware.Logging(),
		middleware.RequestId(),
	)
	// Run the other tests in m.Run()
	os.Exit(m.Run())
}
