package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/rezbow/rssagg/views"
)

func render(ctx context.Context, w http.ResponseWriter, component templ.Component) {
	views.Base(component, "RSS-AGG").Render(ctx, w)
}

func redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func notfound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func serverError(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func extractID(r *http.Request) (int, error) {
	idValue := r.PathValue("id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		return 0, fmt.Errorf("failed to extract id from path: %w", err)
	}
	if id < 0 {
		return 0, errors.New("failed to extrat id form path: invliad id")
	}
	return id, nil
}
