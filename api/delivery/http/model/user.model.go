package model

import (
	"context"

	"github.com/candrairwn/go-pure/api/delivery/http/helper"
	"github.com/google/uuid"
)

type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u UserLoginReq) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)
	if u.Username == "" {
		problems["username"] = "Username is required"
	}

	if u.Password == "" {
		problems["password"] = "Password is required"
	}

	return problems

}

type UserJWT struct {
	Version    int                     `json:"version"`
	Id         uuid.UUID               `json:"id"`
	Username   string                  `json:"username"`
	IdTipeUser string                  `json:"id_tipe_user"`
	IdProdi    helper.Nullable[string] `json:"id_prodi"`
	NamaProdi  helper.Nullable[string] `json:"nama_prodi"`
}
