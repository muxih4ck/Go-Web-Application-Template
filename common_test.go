package main

import (
	"apiserver/config"
	"apiserver/model"
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path.")
)

// This function is used to do setup before executing the test functions
func TestMain(m *testing.M) {

	pflag.Parse()

	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// init db
	model.DB.Init()
	defer model.DB.Close()
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)
	fmt.Println("fewfewf")
	// Run the other tests
	os.Exit(m.Run())
}

// Helper function to create a router during testing
func getRouter(withTemplates bool) *gin.Engine {
	r := gin.Default()
	if withTemplates {
		r.LoadHTMLGlob("templates/*")
		r.Use(setUserStatus())
	}
	return r
}
