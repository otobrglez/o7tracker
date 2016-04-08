package o7tracker

import (
	"appengine"
	"appengine/datastore"
	"net/http"
	"strings"
	"time"
)

// Repository for dealing with DataStore
type Repository struct {
	context appengine.Context
}

// SaveCampaign creates or updates campaign
func (r *Repository) SaveCampaign(campaign *Campaign, id int64) (*datastore.Key, error) {
	campaign.UpdatedAt = time.Now()
	if id == 0 {
		campaign.CreatedAt = time.Now()
	}

	valid, err := campaign.Valid()
	if !valid || err != nil {
		return nil, err
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
func (r *Repository) ListCampaigns(params map[string][]string) ([]Campaign, error) {
	d := datastore.NewQuery("Campaign")

	// Filter per platform
	queryPlatforms := params["platforms"]

	if queryPlatforms != nil && len(queryPlatforms) != 0 {
		pomPlatforms := strings.Split(queryPlatforms[0], ",")
		for _, platform := range pomPlatforms {
			if IsPlatformSupported(platform) {
				d = d.Filter("platforms = ", platform)
			}
		}
	}

	//TODO: Implement pagination
	// d = d.Limit(10)
	d = d.Order("-CreatedAt")

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



// SaveClick
func (r *Repository) SaveClick(platformId int64, r *http.Request) (Click, error) {


	return nil
}
