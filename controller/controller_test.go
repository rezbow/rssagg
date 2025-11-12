package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/rezbow/rssagg/models"
)

type StubFeedRepo struct {
	addFeedCalls  []models.Feed
	editFeedCalls []int
}

func (repo *StubFeedRepo) GetFeeds() []models.Feed {
	return []models.Feed{
		{ID: 1, Title: "New York Times"},
		{ID: 2, Title: "The Guardian"},
		{ID: 3, Title: "Hacker News"},
	}
}

func (repo *StubFeedRepo) GetFeed(int) (*models.Feed, error) {
	return &models.Feed{ID: 4, Title: "The Washingon Post"}, nil
}

func (repo *StubFeedRepo) AddFeed(feed *models.Feed) error {
	feed.ID = 69
	repo.addFeedCalls = append(repo.addFeedCalls, *feed)
	return nil
}

func (repo *StubFeedRepo) EditFeed(feed *models.Feed) error {
	repo.editFeedCalls = append(repo.editFeedCalls, feed.ID)
	return nil
}

func (repo *StubFeedRepo) assertFeedEdited(t *testing.T, feed models.Feed) {
	t.Helper()
	if len(repo.editFeedCalls) != 1 {
		t.Fatalf("got %d calls to EditFeed, wanted %d", len(repo.editFeedCalls), 1)
	}
	if repo.addFeedCalls[0].ID != feed.ID {
		t.Errorf("feed with id %d edited, wanted %d to be edited", repo.editFeedCalls[0], feed.ID)
	}
}

func (repo *StubFeedRepo) assertFeedAdded(t *testing.T, feed models.Feed) {
	t.Helper()
	if len(repo.addFeedCalls) != 1 {
		t.Fatalf("got %d calls to AddFeed, wanted %d", len(repo.addFeedCalls), 1)
	}
	if repo.addFeedCalls[0].Title != feed.Title {
		t.Errorf("feed %v added, wanted %v to be added", repo.addFeedCalls[0], feed)
	}
}

func TestController(t *testing.T) {
	repo := &StubFeedRepo{}
	controller := NewController(repo)
	t.Run("returns list of all feeds", func(t *testing.T) {
		req := newGetRequest("/feeds")
		res := httptest.NewRecorder()
		controller.ServeHTTP(res, req)
		assertCode(t, res.Code, http.StatusOK)
		cupaloy.SnapshotT(t, res.Body.String())
	})

	t.Run("returns the information about a single feed", func(t *testing.T) {
		req := newGetRequest("/feeds/1")
		res := httptest.NewRecorder()
		controller.ServeHTTP(res, req)
		assertCode(t, res.Code, http.StatusOK)
		cupaloy.SnapshotT(t, res.Body.String())
	})

	t.Run("returns feed form", func(t *testing.T) {
		req := newGetRequest("/feeds/new")
		res := httptest.NewRecorder()
		controller.ServeHTTP(res, req)
		assertCode(t, res.Code, http.StatusOK)
		cupaloy.SnapshotT(t, res.Body.String())
	})

	t.Run("add a new feed", func(t *testing.T) {
		newFeed := models.Feed{
			Title: "Reddit",
		}
		req := newFeedRequest(newFeed)
		res := httptest.NewRecorder()
		controller.ServeHTTP(res, req)
		repo.assertFeedAdded(t, newFeed)
		assertRedirect(t, res, "/feeds/69")
	})

	t.Run("returns feed edit form", func(t *testing.T) {
		req := newGetRequest("/feeds/69/edit")
		res := httptest.NewRecorder()
		controller.ServeHTTP(res, req)
		assertCode(t, res.Code, http.StatusOK)
		cupaloy.SnapshotT(t, res.Body.String())
	})

	t.Run("must edit a feed", func(t *testing.T) {
		feed := models.Feed{
			ID:    69,
			Title: "Reddit",
		}
		req := editFeedRequest(feed)
		res := httptest.NewRecorder()
		controller.ServeHTTP(res, req)
		assertRedirect(t, res, "/feeds/69")
		repo.assertFeedEdited(t, feed)
	})
}

func editFeedRequest(feed models.Feed) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/feeds/%d/edit", feed.ID), feedToFormReader(feed))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req

}

func newFeedRequest(feed models.Feed) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/feeds", feedToFormReader(feed))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func feedToFormReader(feed models.Feed) *strings.Reader {
	form := url.Values{}
	form.Add("title", feed.Title)
	return strings.NewReader(form.Encode())
}

func newGetRequest(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return req
}

func assertRedirect(t testing.TB, res *httptest.ResponseRecorder, location string) {
	t.Helper()
	assertCode(t, res.Code, http.StatusSeeOther)
	locationHeader := res.Header().Get("Location")
	if locationHeader != location {
		t.Errorf("got location header %q, wanted %q", locationHeader, location)
	}
}

func assertCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got code %d, want %d", got, want)
	}
}
