package o7tracker

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"fmt"
	"net/http"
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

// Repository for dealing with DataStore
type Repository struct {
	context appengine.Context
}

// IsPlatformSupported returns true if platform is supported by the system
func IsPlatformSupported(platform string) bool {
	i := sort.SearchStrings(SupportedPlatforms, platform)
	return i < len(SupportedPlatforms) && SupportedPlatforms[i] == platform
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

// SaveCampaign creates or updates campaign
func (r *Repository) SaveCampaign(campaign *Campaign, id int64) (*datastore.Key, error) {
	campaign.UpdatedAt = time.Now()
	if id == 0 {
		campaign.CreatedAt = time.Now()
	}

	valid, error := campaign.Valid()
	if !valid || error != nil {
		return nil, error
	}

	key := datastore.NewKey(r.context, "Campaign", "", id, nil)
	newKey, err := datastore.Put(r.context, key, campaign)
	if err != nil {
		return nil, err
	}

	campaign.ID = newKey.IntID()
	return newKey, nil
}

// ListCampaigns lists all campaigns
func (r *Repository) ListCampaigns() ([]Campaign, error) {
	d := datastore.NewQuery("Campaign").Limit(10).Order("-CreatedAt")
	campaigns := make([]Campaign, 0, 10)
	keys, err := d.GetAll(r.context, &campaigns)
	if err != nil {
		return campaigns, err
	}

	for i, key := range keys {
		campaigns[i].ID = key.IntID()
	}

	return campaigns, nil
}

// GetCampaign returns single campaign
func (r *Repository) GetCampaign(id int64) (Campaign, error) {
	key := datastore.NewKey(r.context, "Campaign", "", id, nil)
	campaign := Campaign{}
	err := datastore.Get(r.context, key, &campaign)
	return campaign, err
}

// DeleteCampaign deletes campaign
func (r *Repository) DeleteCampaign(id int64) error {
	key := datastore.NewKey(r.context, "Campaign", "", id, nil)
	r.context.Infof("Deleting %s", key.String())
	return datastore.Delete(r.context, key)
}


// AddClickFrom adds click
func (r *Repository) AddClickFrom(request *http.Request) {

}
