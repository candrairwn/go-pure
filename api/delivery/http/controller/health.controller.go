package controller

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/candrairwn/go-pure/api/delivery/http/model"
	"github.com/candrairwn/go-pure/api/utils"
)

type HealthController struct{}

func (c *HealthController) HandleGetHealth(version string) http.HandlerFunc {
	res := model.HealthRes{Version: version}
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

	up := time.Now()

	return func(w http.ResponseWriter, r *http.Request) {

		res.Uptime = time.Since(up).String()
		utils.Encode(w, r, http.StatusOK, res)
	}
}
