package o7tracker

import (
  "appengine"
  "appengine/datastore"
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
  campaign := NewCampaign()
  err := datastore.Get(r.context, key, &campaign)
  campaign.ID = key.IntID()
  return campaign, err
}

// GetCampaignWithDBStats retrieves campaign with stats from Datastore
func (r *Repository) GetCampaignComputeCounts(id int64) (Campaign, error) {
  key := datastore.NewKey(r.context, "Campaign", "", id, nil)
  campaign := NewCampaign()
  err := datastore.Get(r.context, key, &campaign)
  if err != nil {
    return campaign, err
  }

  basicQuery := datastore.NewQuery("Click").Ancestor(key).KeysOnly()
  count, err := basicQuery.Count(r.context)
  if err != nil {
    return campaign, err
  }

  androidCount, err := basicQuery.
  Filter("Platform =", "Android").
  Count(r.context)
  if err != nil {
    return campaign, err
  }

  iphoneCount, err := basicQuery.
  Filter("Platform =", "IPhone").
  Count(r.context)
  if err != nil {
    return campaign, err
  }

  windowsPhoneCount, err := basicQuery.
  Filter("Platform =", "WindowsPhone").
  Count(r.context)
  if err != nil {
    return campaign, err
  }

  campaign.ID = key.IntID()
  campaign.ClickCount = count
  campaign.AndroidClickCount = androidCount
  campaign.IPhoneClickCount = iphoneCount
  campaign.WindowsPhoneClickCount = windowsPhoneCount

  return campaign, err
}

// DeleteCampaign deletes campaign
func (r *Repository) DeleteCampaign(id int64) error {
  key := datastore.NewKey(r.context, "Campaign", "", id, nil)
  r.context.Infof("Campaign %d was successfuly deleted.", key.IntID())
  return datastore.Delete(r.context, key)
}

// TrackClick stores click
func (r *Repository) TrackClick(campaignID int64, click *Click) (Campaign, error) {
  parentKey := datastore.NewKey(r.context, "Campaign", "", campaignID, nil)
  key := datastore.NewIncompleteKey(r.context, "Click", parentKey)

  campaign, err := r.GetCampaign(campaignID)
  if err != nil {
    return campaign, err
  }

  outKey, err := datastore.Put(r.context, key, click)
  if err != nil {
    return campaign, err
  }

  r.context.Infof("outKey %s", outKey.String())

  return campaign, nil
}
