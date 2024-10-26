package model

type HealthRes struct {
	Version    string `json:"version"`
	Uptime     string `json:"uptime"`
	DirtyBuild bool   `json:"dirty_build"`
	DbStatus   string `json:"db_status"`
}
