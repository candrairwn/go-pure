package controller

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/candrairwn/go-pure/api/delivery/http/model"
	"github.com/candrairwn/go-pure/api/utils"
	"gorm.io/gorm"
)

type HealthController struct {
	Db      *gorm.DB
	Version string
}

func NewHealthController(db *gorm.DB, version string) *HealthController {
	return &HealthController{Db: db, Version: version}
}

func (c *HealthController) HandleGetHealth() http.HandlerFunc {
	res := model.HealthRes{Version: c.Version}
	buildInfo, _ := debug.ReadBuildInfo()
	for _, kv := range buildInfo.Settings {
		if kv.Value == "" {
			continue
		}
		switch kv.Key {
		case "vcs.modified":
			res.DirtyBuild = kv.Value == "true"
		}
	}

	res.DbStatus = "Up"

	sqlDb, err := c.Db.DB()
	if err != nil {
		res.DbStatus = "Error"
	} else {
		if err := utils.CheckDBStatus(sqlDb); err != nil {
			res.DbStatus = "Down"
		}
	}

	up := time.Now()

	return func(w http.ResponseWriter, r *http.Request) {
		res.Uptime = time.Since(up).String()
		utils.EncodeWithWrapper(w, r, http.StatusOK, res, nil)
	}
}
