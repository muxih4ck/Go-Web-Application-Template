package main

import (
	"apiserver/config"
	"apiserver/handler/user"
	"apiserver/model"
	"apiserver/router"
	"apiserver/router/middleware"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/viper"

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

func TestLogin(t *testing.T) {
	uri := "/login"
	u := user.CreateRequest{
		Username: "admin",
		Password: "admin",
	}
	jsonByte, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Test Error: %s", err.Error())
	}
	req := httptest.NewRequest(http.MethodPost, uri, bytes.NewReader(jsonByte))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	g.ServeHTTP(w, req)
	result := w.Result()

	defer result.Body.Close()

	// 读取响应body,获取token
	var data model.LoginResponse
	bodyByte, _ := ioutil.ReadAll(result.Body)
	if err := json.Unmarshal(bodyByte, &data); err != nil {
		t.Errorf("Test error: Get LoginResponse Error:%s", err.Error())
	}
	token = data.Data.Token

	if result.StatusCode != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", result.StatusCode)
	}
}

func TestCreate(t *testing.T) {
	uri := "/v1/user"

	username = strconv.FormatInt(time.Now().Unix(), 10)
	password = strconv.FormatInt(time.Now().Unix(), 10)

	u := user.CreateRequest{
		Username: username,
		Password: password,
	}
	jsonByte, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Test Error: %s", err.Error())
	}
	req := httptest.NewRequest(http.MethodPost, uri, bytes.NewReader(jsonByte))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	g.ServeHTTP(w, req)
	result := w.Result()

	// GetUid
	user, err := model.GetUser(username)
	if err != nil {
		t.Errorf("Test Error: %s", err.Error())
	}
	uid = user.Id

	if result.StatusCode != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", result.StatusCode)
	}
}

func TestGet(t *testing.T) {
	uri := "/v1/user/" + username

	req := httptest.NewRequest(http.MethodGet, uri, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	g.ServeHTTP(w, req)
	result := w.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", result.StatusCode)
	}
}

func TestUpdate(t *testing.T) {
	uri := "/v1/user/" + strconv.FormatInt(int64(uid), 10)
	u := user.CreateRequest{
		Username: "test" + username,
		Password: "test" + password,
	}
	jsonByte, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Test Error: %s", err.Error())
	}
	req := httptest.NewRequest(http.MethodPut, uri, bytes.NewReader(jsonByte))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	g.ServeHTTP(w, req)
	result := w.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", result.StatusCode)
	}
}

func TestList(t *testing.T) {
	uri := "/v1/user"

	req := httptest.NewRequest(http.MethodGet, uri, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	g.ServeHTTP(w, req)
	result := w.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", result.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	uri := "/v1/user/" + strconv.FormatInt(int64(uid), 10)

	req := httptest.NewRequest(http.MethodDelete, uri, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	g.ServeHTTP(w, req)
	result := w.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", result.StatusCode)
	}

}
