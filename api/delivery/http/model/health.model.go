package model

import "time"

type HealthRes struct {
	Version        string    `json:"version"`
	Uptime         string    `json:"uptime"`
	LastCommitHash string    `json:"last_commit_hash"`
	LastCommitTime time.Time `json:"last_commit_time"`
	DirtyBuild     bool      `json:"dirty_build"`
}
