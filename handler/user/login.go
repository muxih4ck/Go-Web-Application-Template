package user

import (
	. "github.com/muxih4ck/Go-Web-Application-Template/handler"
	"github.com/muxih4ck/Go-Web-Application-Template/model"
	"github.com/muxih4ck/Go-Web-Application-Template/pkg/auth"
	"github.com/muxih4ck/Go-Web-Application-Template/pkg/errno"
	"github.com/muxih4ck/Go-Web-Application-Template/pkg/token"

	"github.com/gin-gonic/gin"
)

// Login generates the authentication token
// if the password was matched with the specified account.
func Login(c *gin.Context) {
	// Binding the data with the user struct.
	var u model.UserModel
	if err := c.Bind(&u); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	// Get the user information by the login username.
	d, err := model.GetUser(u.Username)
	if err != nil {
		SendError(c, errno.ErrUserNotFound, nil, err.Error())
		return
	}

	// Compare the login password with the user password.
	// 业务逻辑异常，使用 SendResponse 发送 200 请求 + 自定义错误码
	if err := auth.Compare(d.Password, u.Password); err != nil {
		SendResponse(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	// Sign the json web token.
	t, err := token.Sign(c, token.Context{ID: d.Id, Username: d.Username}, "")
	if err != nil {
		SendError(c, errno.ErrToken, nil, err.Error())
		return
	}

	SendResponse(c, nil, model.Token{Token: t})
}
