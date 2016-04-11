package o7tracker

import (
	"appengine/aetest"
	"testing"
)

func TestTrack(t *testing.T) {
	context, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}

	r := Repository{context}

	campaign := Campaign{
		RedirectURL: "https://github.com",
		Platforms:   SupportedPlatforms,
	}

	// Save campaign
	key, err := r.SaveCampaign(&campaign, 0)
	if err != nil {
		t.Fatal(err)
	}

	// Track campaign click
	if _, err := r.Track(&Click{
		CampaignID: key.IntID(),
		Platform:   SupportedPlatforms[0],
	}); err != nil {
		t.Fatal(err)
	}

	// Get stats for campaign
	stats, err := r.GetCampaignComputeCounts(key.IntID())
	if err != nil {
		t.Fatal(err)
	}

	if stats.ClickCount != 1 {
		t.Fatal("Stats ware not saved.")
	}

	// List campaign
	campaigns, err := r.ListCampaigns(map[string][]string{})
	if err != nil {
		t.Fatal(err)
	}

	if len(campaigns) != 1 {
		t.Fatalf("Campaigns %d", len(campaigns))
	}

	// Delete campaign
	if err := r.DeleteCampaign(key.IntID()); err != nil {
		t.Fatal(err)
	}
}
