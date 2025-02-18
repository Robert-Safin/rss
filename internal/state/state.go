package state

import (
	"rss/internal/config"
	"rss/internal/database"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}
