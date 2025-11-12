package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rezbow/rssagg/views/feeds"
)

type Controller struct {
	repo Repo
	http.Handler
}

func NewController(repo Repo) *Controller {
	controller := &Controller{
		repo: repo,
	}

	router := http.NewServeMux()

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))
	router.Handle("GET /feeds", http.HandlerFunc(controller.feeds))
	router.Handle("POST /feeds", http.HandlerFunc(controller.newFeed))
	router.Handle("GET /feeds/new", http.HandlerFunc(controller.feedForm))
	router.Handle("GET /feeds/{id}", http.HandlerFunc(controller.feed))
	router.Handle("GET /feeds/{id}/edit", http.HandlerFunc(controller.feedEdit))
	router.Handle("POST /feeds/{id}/edit", http.HandlerFunc(controller.editFeed))

	controller.Handler = router

	return controller
}

func (c *Controller) newFeed(w http.ResponseWriter, r *http.Request) {
	form := feeds.FormFromRequest(r)
	if !form.Valid() {
		render(r.Context(), w, feeds.Form(form))
		return
	}
	feed := form.ToFeed()
	if err := c.repo.AddFeed(feed); err != nil {
		log.Println(err)
		serverError(w)
		return
	}
	redirect(w, r, fmt.Sprintf("/feeds/%d", feed.ID))
}

func (c *Controller) feedForm(w http.ResponseWriter, r *http.Request) {
	render(r.Context(), w, feeds.Form(feeds.NewFeedForm()))
}

func (c *Controller) feeds(w http.ResponseWriter, r *http.Request) {
	render(r.Context(), w, feeds.Feeds(c.repo.GetFeeds()))
}

func (c *Controller) feed(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r)
	if err != nil {
		notfound(w, r)
		return
	}
	feed, err := c.repo.GetFeed(id)
	if err != nil || feed == nil {
		notfound(w, r)
		return
	}
	render(r.Context(), w, feeds.Feed(feed))
}

func (c *Controller) feedEdit(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r)
	if err != nil {
		notfound(w, r)
		return
	}
	feed, err := c.repo.GetFeed(id)
	if err != nil || feed == nil {
		notfound(w, r)
		return
	}
	render(r.Context(), w, feeds.EditForm(feeds.NewFormFromFeed(feed)))
}

func (c *Controller) editFeed(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r)
	if err != nil {
		notfound(w, r)
		return
	}
	form := feeds.FormFromRequest(r)
	form.SetId(id)
	if !form.Valid() {
		render(r.Context(), w, feeds.EditForm(form))
		return
	}
	feed := form.ToFeed()
	feed.ID = id
	if err := c.repo.EditFeed(feed); err != nil {
		notfound(w, r)
		return
	}
	redirect(w, r, fmt.Sprintf("/feeds/%d", id))
}
