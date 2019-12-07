package user

import (
	. "github.com/muxih4ck/Go-Web-Application-Template/handler"
	"github.com/muxih4ck/Go-Web-Application-Template/log"
	"github.com/muxih4ck/Go-Web-Application-Template/model"
	"github.com/muxih4ck/Go-Web-Application-Template/pkg/errno"
	"github.com/muxih4ck/Go-Web-Application-Template/util"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// Create creates a new user account.
func Create(c *gin.Context) {
	log.Info("User Create function called.",
		zap.String("X-Request-Id", util.GetReqID(c)))
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	u := model.UserModel{
		Username: r.Username,
		Password: r.Password,
	}

	// Validate the data.
	if err := u.Validate(); err != nil {
		SendError(c, errno.ErrValidation, nil, err.Error())
		return
	}

	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		SendError(c, errno.ErrEncrypt, nil, err.Error())
		return
	}
	// Insert the user to the database.
	if err := u.Create(); err != nil {
		SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	rsp := CreateResponse{
		Username: r.Username,
	}

	// Show the user information.
	SendResponse(c, nil, rsp)
}
