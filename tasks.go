package o7tracker

import (
	"appengine"
	"net/http"
)

// UpdateStats method that handles HTTP request for stats update
func UpdateStats(w http.ResponseWriter, r *http.Request) {
	repository := Repository{appengine.NewContext(r)}
	if err := repository.UpdateStats(); err != nil {
		ErrorToJSON(w, err)
		return
	}

	w.Write([]byte("Updated."))
}
