package o7tracker

import (
	"appengine"
	"net/http"
  "strconv"
)

func init() {
	http.HandleFunc("/track/", handleTrack)
	http.HandleFunc("/admin/campaigns", AdminCampaigns)
	http.HandleFunc("/admin/campaigns/", AdminCampaigns)
}

func handleTrack(w http.ResponseWriter, r *http.Request) {
  id, err := strconv.Atoi(r.URL.Path[(len("/admin/track") + 1):])
  if err != nil {
    ErrorToJSON(w, err)
    return
  }

  repository := Repository{appengine.NewContext(r)}
  click, err := repository.SaveClick(int64(id), r)
  if err != nil {
    ErrorToJSON(w, err)
    return
  }


}
