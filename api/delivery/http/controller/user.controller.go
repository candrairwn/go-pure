package controller

import (
	"net/http"

	"github.com/candrairwn/go-pure/api/delivery/http/helper"
	"github.com/candrairwn/go-pure/api/delivery/http/model"
	"github.com/candrairwn/go-pure/api/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserController struct {
	Log *zap.SugaredLogger
}

func NewUserController(log *zap.SugaredLogger) *UserController {
	return &UserController{
		Log: log,
	}
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {

	userLoginReq, problems, err := utils.DecodeValid[model.UserLoginReq](r)
	if err != nil {
		utils.EncodeWithWrapper(w, r, http.StatusBadRequest, model.UserLoginReq{}, map[string]interface{}{
			"error":    err.Error(),
			"problems": problems,
		})
		return
	}

	token, err := utils.CreateAccessToken(model.UserJWT{
		Version:    1,
		Id:         uuid.New(),
		Username:   userLoginReq.Username,
		IdTipeUser: "SU",
		IdProdi:    helper.NewValue("", false),
		NamaProdi:  helper.NewValue("", false),
	}, 3600, c.Log)

	if err != nil {
		utils.EncodeWithWrapper(w, r, http.StatusInternalServerError, model.UserLoginReq{}, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	utils.EncodeWithWrapper(w, r, http.StatusOK, model.UserLoginReq{
		Username: userLoginReq.Username,
		Password: userLoginReq.Password,
	}, map[string]interface{}{
		"access_token": token,
	})
}
