package o7tracker

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"time"
)

// SupportedPlatforms contains list of all currently supported platforms
var SupportedPlatforms = []string{"Android", "IPhone", "WindowsPhone"}

// Campaign that holds "campaign" information.
type Campaign struct {
	ID             int64     `json:"id"`
	RedirectURL    string    `datastore:"" json:"redirect_url"`
	NumberOfClicks int       `datastore:",noindex" json:"number_of_clicks"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Platforms      []string  `datastore:"platforms" json:"platforms" datastore_type:"StringList" verbose_name:"platforms"`
}

// Valid validates Campaign class
func (c *Campaign) Valid() (bool, error) {
	// Check if all platforms are supported.
	areAllSupported := true
	for _, platform := range c.Platforms {
		if !IsPlatformSupported(platform) {
			areAllSupported = false
		}
	}

	if areAllSupported != true {
		return false, fmt.Errorf("Submited platforms are invalid %s", c.Platforms)
	}

	// Check if URL is present
	if c.RedirectURL == "" {
		return false, errors.New(`"redirect_url" is required parameter`)
	}

	// Check if URL is valid
	_, err := url.ParseRequestURI(c.RedirectURL)
	if err != nil {
		return false, errors.New(`Invalid "redirect_url"`)
	}

	return true, nil
}

// IsPlatformSupported returns true if platform is supported by the system
func IsPlatformSupported(platform string) bool {
	i := sort.SearchStrings(SupportedPlatforms, platform)
	return i < len(SupportedPlatforms) && SupportedPlatforms[i] == platform
}
