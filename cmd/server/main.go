package main

import (
	"net/http"

	"github.com/rezbow/rssagg/controller"
	inmemory "github.com/rezbow/rssagg/repo/in_memory"
)

func main() {
	repo := inmemory.NewInMemoryRepo()
	http.ListenAndServe(":8080", controller.NewController(repo))
}
