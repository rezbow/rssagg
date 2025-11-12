package inmemory

import (
	"errors"

	"github.com/rezbow/rssagg/models"
)

type InMemoryRepo struct {
	feeds []models.Feed
}

func (r *InMemoryRepo) AddFeed(feed *models.Feed) error {
	feed.ID = len(r.feeds) + 1
	r.feeds = append(r.feeds, *feed)
	return nil
}

func (r *InMemoryRepo) GetFeeds() []models.Feed {
	return r.feeds
}

func (r *InMemoryRepo) GetFeed(id int) (*models.Feed, error) {
	if !r.exists(id) {
		return nil, errors.New("feed not found")
	}
	return &r.feeds[id-1], nil
}

func (r *InMemoryRepo) EditFeed(feed *models.Feed) error {
	if !r.exists(feed.ID) {
		return errors.New("feed not found")
	}
	f := &r.feeds[feed.ID-1]
	f.Title = feed.Title
	return nil
}

func (r *InMemoryRepo) exists(id int) bool {
	return id >= 1 && id <= len(r.feeds)
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		feeds: []models.Feed{
			{ID: 1, Title: "New York Times"},
			{ID: 2, Title: "The Guardian"},
			{ID: 3, Title: "Hackernews"},
		},
	}
}
