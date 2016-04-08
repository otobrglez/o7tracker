package o7tracker

import (
	"net/http"
)

func init() {
	http.HandleFunc("/track", handleTrack)
	http.HandleFunc("/admin/campaigns", AdminCampaigns)
	http.HandleFunc("/admin/campaigns/", AdminCampaigns)
}

func handleTrack(w http.ResponseWriter, r *http.Request) {

}
