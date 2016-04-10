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

const (
	defaultFallbackURL    = "http://outfit7.com"
	adminTracksPath       = "/track/"
	adminTracksPathLength = len(adminTracksPath)
)

// ClickTrack controller that handles all the HTTP magic.
func ClickTrack(w http.ResponseWriter, r *http.Request) {
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
	campaign, err := repository.Track(&Click{
		CampaignID: id,
		Platform:   platform,
		CreatedAt:  time.Now(),
	})
	if err != nil {
		handleError(context, w, r, err)
		return
	}

	info := fmt.Sprintf("platform_id:%d platform:%s redirect_url:%s", id, platform, campaign.RedirectURL)
	context.Infof(info)
	http.Redirect(w, r, campaign.RedirectURL, http.StatusTemporaryRedirect)
}

func handleError(context appengine.Context, w http.ResponseWriter, r *http.Request, err error) {
	fallbackURL := os.Getenv("FALLBACK_URL")
	if fallbackURL == "" {
		fallbackURL = defaultFallbackURL
	}

	context.Infof("Falling back to %s ~> %s", r.RequestURI, fallbackURL)
	http.Redirect(w, r, fallbackURL, http.StatusTemporaryRedirect)
}
