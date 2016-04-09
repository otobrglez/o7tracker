package o7tracker

import (
  "appengine"
  "errors"
  "fmt"
  "net/http"
  "os"
  "strconv"
  "time"
)

func init() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    context := appengine.NewContext(r)
    handleError(context, w, r, errors.New(`Failing on "*" path.`))
  })

  http.HandleFunc("/track/", handleTrack)
  http.HandleFunc("/admin/campaigns", AdminCampaigns)
  http.HandleFunc("/admin/campaigns/", AdminCampaigns)
}

func handleError(context appengine.Context, w http.ResponseWriter, r *http.Request, err error) {
  defaultFallbackURL := "http://outfit7.com"
  fallbackURL := os.Getenv("FALLBACK_URL")
  if fallbackURL == "" {
    fallbackURL = defaultFallbackURL
  }

  context.Infof("Falling back to %s ~> %s", r.RequestURI, fallbackURL)
  http.Redirect(w, r, fallbackURL, http.StatusTemporaryRedirect)
}

const adminTracksPath = "/track/"
const adminTracksPathLength = len(adminTracksPath)

func handleTrack(w http.ResponseWriter, r *http.Request) {
  context := appengine.NewContext(r)

  id, err := strconv.ParseInt(r.URL.Path[adminTracksPathLength:], 10, 64)
  if err != nil {
    handleError(context, w, r, err)
    return
  }

  platform := r.URL.Query().Get("platform")
  if platform == "" {
    handleError(context, w, r, errors.New("Missing platform id."))
    return
  }

  repository := Repository{context}
  campaign, err := repository.TrackClick(id, &Click{
    Platform:  platform,
    CreatedAt: time.Now(),
  })
  if err != nil {
    handleError(context, w, r, err)
    return
  }

  info := fmt.Sprintf("platform_id:%d platform:%s", id, platform)
  context.Infof(info)
  http.Redirect(w,r,campaign.RedirectURL, http.StatusTemporaryRedirect)
}
