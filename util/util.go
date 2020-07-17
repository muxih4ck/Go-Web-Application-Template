package util

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
)

func GenShortId() (string, error) {
	return shortid.Generate()
}

func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestID, ok := v.(string); ok {
		return requestID
	}
	return ""
}

func GetProjectAbsPath() string {
	var (
		path string
		err  error
	)
	defer func() {
		if err != nil {
			panic(fmt.Sprintf("Find config file by using functiong GetProjectAbsPath error :%+v", err))
		}
	}()
	path, err = os.Getwd()
	fmt.Println(path)
	return path
}
