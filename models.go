package o7tracker

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
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
	ID          int64     `json:"id"`
	RedirectURL string    `datastore:"" datastore_type:"Link" json:"redirect_url" `
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Platforms   []string  `datastore:"platforms" json:"platforms" datastore_type:"StringList" verbose_name:"platforms"`

	ClickCount             int `datastore:",noindex" json:"click_count"`
	AndroidClickCount      int `datastore:",noindex" json:"android_click_count"`
	IPhoneClickCount       int `datastore:",noindex" json:"iphone_click_count"`
	WindowsPhoneClickCount int `datastore:",noindex" json:"windowsphone_click_count"`
}

// Click model
type Click struct {
	CampaignID int64     `datastore:"-" json:"-"`
	Platform   string    `datastore:"" json:"platform"`
	UserAgent  string    `datastore:"" json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

// NewCampaign creates new empty campaign
func NewCampaign() Campaign {
	campaign := Campaign{
		ClickCount:             0,
		AndroidClickCount:      0,
		IPhoneClickCount:       0,
		WindowsPhoneClickCount: 0,
	}
	return campaign
}

// Valid validates Campaign instance
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

// CacheKeys returns keys for memcache
func (c *Campaign) CacheKeys() (keys []string) {
	for _, platform := range c.Platforms {
		keys = append(keys, fmt.Sprintf("%d-%s", c.ID, platform))
	}
	return
}

// CacheItems returns collection of items to be stored to memcache
func (c *Campaign) CacheItems() []*memcache.Item {
	var items []*memcache.Item
	for _, platform := range c.Platforms {
		click := Click{
			CampaignID: c.ID,
			Platform:   platform,
		}

		item := &memcache.Item{
			Key:   click.CacheKey(),
			Value: []byte(c.RedirectURL),
		}

		items = append(items, item)
	}

	return items
}

// Valid validates Click instance
func (c *Click) Valid() (bool, error) {
	if c.Platform == "" {
		return false, errors.New("Missing platform")
	}

	if !IsPlatformSupported(c.Platform) {
		return false, fmt.Errorf("Platform %s is not supported", c.Platform)
	}

	if c.CampaignID == 0 {
		return false, errors.New("Missing CampaignID")
	}

	return true, nil
}

// DatastoreKey returns Click key for datastore
func (c *Click) DatastoreKey(context appengine.Context) *datastore.Key {
	parentKey := datastore.NewKey(context, "Campaign", "", c.CampaignID, nil)
	return datastore.NewIncompleteKey(context, "Click", parentKey)
}

// CacheKey returns Click key for memcache
func (c *Click) CacheKey() string {
	return fmt.Sprintf("%d-%s", c.CampaignID, c.Platform)
}

// IsPlatformSupported returns true if platform is supported by the system
func IsPlatformSupported(platform string) bool {
	i := sort.SearchStrings(SupportedPlatforms, platform)
	return i < len(SupportedPlatforms) && SupportedPlatforms[i] == platform
}
