package controller

import "github.com/rezbow/rssagg/models"

type Repo interface {
	GetFeeds() []models.Feed
	GetFeed(int) (*models.Feed, error)
	AddFeed(*models.Feed) error
	EditFeed(*models.Feed) error
}
