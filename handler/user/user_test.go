package user

import (
	"apiserver/config"
	"apiserver/model"
	"apiserver/router/middleware"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	g           *gin.Engine
	tokenString string
	username    string
	password    string
	uid         uint64
)

func TestMain(m *testing.M) {

	// init config
	if err := config.Init(""); err != nil {
		panic(err)
	}
	// init db
	model.DB.Init()
	defer model.DB.Close()

	os.Exit(m.Run())
}
func TestLogin(t *testing.T) {
	g := getRouter(true)

	uri := "/login"
	u := CreateRequest{
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
	// result := w.Result()

	// defer result.Body.Close()

	// 读取响应body,获取tokenString
	var data model.LoginResponse
	// bodyByte, _ := ioutil.ReadAll(result.Body)
	fmt.Println("fwefwef", w.Body.String())
	if err := json.Unmarshal([]byte(w.Body.String()), &data); err != nil {
		t.Errorf("Test error: Get LoginResponse Error:%s", err.Error())
	}
	tokenString = data.Data.Token

	if w.Code != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", w.Code)
	}
}

func TestCreate(t *testing.T) {
	g := getRouter(true)
	uri := "/v1/user"

	username = strconv.FormatInt(time.Now().Unix(), 10)
	password = strconv.FormatInt(time.Now().Unix(), 10)

	u := CreateRequest{
		Username: username,
		Password: password,
	}
	jsonByte, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Test Error: %s", err.Error())
	}
	req := httptest.NewRequest(http.MethodPost, uri, bytes.NewReader(jsonByte))
	req.Header.Set("Authorization", "Bearer "+tokenString)
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
	g := getRouter(true)
	uri := "/v1/user/" + username

	req := httptest.NewRequest(http.MethodGet, uri, nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	g.ServeHTTP(w, req)
	result := w.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", result.StatusCode)
	}
}

func TestUpdate(t *testing.T) {
	g := getRouter(true)
	uri := "/v1/user/" + strconv.FormatInt(int64(uid), 10)
	u := CreateRequest{
		Username: "test" + username,
		Password: "test" + password,
	}
	jsonByte, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Test Error: %s", err.Error())
	}
	req := httptest.NewRequest(http.MethodPut, uri, bytes.NewReader(jsonByte))
	req.Header.Set("Authorization", "Bearer "+tokenString)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	g.ServeHTTP(w, req)
	result := w.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", result.StatusCode)
	}
}

func TestList(t *testing.T) {
	g := getRouter(true)
	uri := "/v1/user"

	req := httptest.NewRequest(http.MethodGet, uri, nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	g.ServeHTTP(w, req)
	result := w.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", result.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	g := getRouter(true)
	uri := "/v1/user/" + strconv.FormatInt(int64(uid), 10)

	req := httptest.NewRequest(http.MethodDelete, uri, nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	g.ServeHTTP(w, req)
	result := w.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("Test Error: StatusCode Error:%d", result.StatusCode)
	}

}

// Helper function to create a router during testing
func getRouter(withRouter bool) *gin.Engine {
	g = gin.New()
	if withRouter {
		loadRouters(
			// Cores.
			g,

			// Middlwares.
			middleware.Logging(),
			middleware.RequestId(),
		)
	}
	return g
}

// Load loads the middlewares, routes, handlers about Test
func loadRouters(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// api for authentication functionalities
	g.POST("/login", Login)

	// The user handlers, requiring authentication
	u := g.Group("/v1/user")
	u.Use(middleware.AuthMiddleware())
	{
		u.POST("", Create)
		u.DELETE("/:id", Delete)
		u.PUT("/:id", Update)
		u.GET("", List)
		u.GET("/:username", Get)
	}

	return g
}
