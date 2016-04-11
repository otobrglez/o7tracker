package o7tracker

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
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

	// Save to Datastore
	newKey, err := datastore.Put(r.context, key, campaign)
	if err != nil {
		return nil, err
	}

	// Delete & Save from memcache
	if id != 0 {
		if err := memcache.DeleteMulti(r.context, campaign.CacheKeys()); err != nil {
			r.context.Infof("Can't delete multi keys for %d", campaign.ID)
		}
	}

	if err := memcache.SetMulti(r.context, campaign.CacheItems()); err != nil {
		return nil, err
	}

	campaign.ID = newKey.IntID()
	r.context.Infof("Campaign %d was successfuly saved.", campaign.ID)
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

	var campaigns []Campaign
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
	campaign := NewCampaign()
	err := datastore.Get(r.context, key, &campaign)
	campaign.ID = key.IntID()
	return campaign, err
}

// GetCampaignComputeCounts retrieves campaign with stats from Datastore
func (r *Repository) GetCampaignComputeCounts(id int64) (Campaign, error) {
	key := datastore.NewKey(r.context, "Campaign", "", id, nil)
	campaign := NewCampaign()
	campaign.ID = key.IntID()

	basicQuery := datastore.NewQuery("Click").Ancestor(key).KeysOnly()
	count, err := basicQuery.Count(r.context)
	if err != nil {
		return campaign, err
	}
	campaign.ClickCount = count

	countA, err := basicQuery.
		Filter("Platform =", "Android").
		Count(r.context)
	if err != nil {
		return campaign, err
	}
	campaign.AndroidClickCount = countA

	countP, err := basicQuery.
		Filter("Platform =", "IPhone").
		Count(r.context)
	if err != nil {
		return campaign, err
	}
	campaign.IPhoneClickCount = countP

	countW, err := basicQuery.
		Filter("Platform =", "WindowsPhone").
		Count(r.context)
	if err != nil {
		return campaign, err
	}
	campaign.WindowsPhoneClickCount = countW

	return campaign, nil
}

// DeleteCampaign deletes campaign
func (r *Repository) DeleteCampaign(id int64) error {
	campaign, err := r.GetCampaign(id)
	if err != nil {
		return err
	}

	// Delete from memcache
	if err := memcache.DeleteMulti(r.context, campaign.CacheKeys()); err != nil {
		r.context.Infof("Can't delete multi keys for %d", campaign.ID)
	}

	// Delete from Datastore
	key := datastore.NewKey(r.context, "Campaign", "", campaign.ID, nil)
	return datastore.Delete(r.context, key)
}

// Track tracks click and returns Campaign with redirection_url or error.
func (r *Repository) Track(click *Click) (Campaign, error) {
	campaign := Campaign{ID: click.CampaignID}

	// Validate click
	valid, err := click.Valid()
	if !valid && err != nil {
		return campaign, err
	}

	// Get from memcache
	inCache := true
	if item, err := memcache.Get(r.context, click.CacheKey()); err == memcache.ErrCacheMiss {
		r.context.Infof("Memcache CacheMiss for %s", click.CacheKey())
		inCache = false
	} else if err != nil {
		return Campaign{}, err
	} else {
		r.context.Infof("Memcache HIT.")
		campaign.RedirectURL = string(item.Value)
	}

	// If not in memcache try to get from Datastore
	if !inCache {
		campaignFromDB, err := r.GetCampaign(campaign.ID)
		if err != nil {
			return campaign, err
		}

		campaign = campaignFromDB

		if err := memcache.SetMulti(r.context, campaign.CacheItems()); err != nil {
			return campaign, err
		}
	}

	// Save click to Datastore
	if _, err := datastore.Put(r.context, click.DatastoreKey(r.context), click); err != nil {
		return campaign, err
	}

	return campaign, nil
}

// UpdateStats periodically updates stats
func (r *Repository) UpdateStats() error {
	campaigns, err := r.ListCampaigns(map[string][]string{})
	if err != nil {
		return err
	}

	for _, campaign := range campaigns {
		if err := r.UpdateCampaignStats(&campaign); err != nil {
			return err
		}
	}

	r.context.Infof("Stats updated.")
	return nil
}

// UpdateCampaignStats updates stats for given campaign and saves them.
func (r *Repository) UpdateCampaignStats(campaign *Campaign) error {
	stats, err := r.GetCampaignComputeCounts(campaign.ID)
	if err != nil {
		return err
	}

	campaign.ClickCount = stats.ClickCount
	campaign.AndroidClickCount = stats.AndroidClickCount
	campaign.IPhoneClickCount = stats.IPhoneClickCount
	campaign.WindowsPhoneClickCount = stats.WindowsPhoneClickCount

	if _, err := r.SaveCampaign(campaign, campaign.ID); err != nil {
		return err
	}

	return nil
}
