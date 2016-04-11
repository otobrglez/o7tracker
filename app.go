package o7tracker

import (
	"appengine"
	"errors"
	"net/http"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleError(appengine.NewContext(r), w, r, errors.New(`Failing on "/" path.`))
	})

	http.HandleFunc("/track/", ClickTrack)
	http.HandleFunc("/admin/campaigns", AdminCampaigns)
	http.HandleFunc("/admin/campaigns/", AdminCampaigns)
	http.HandleFunc("/tasks/update_stats", UpdateStats)
}
