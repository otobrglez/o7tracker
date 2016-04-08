package o7tracker

import (
	"appengine"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// AdminCampaigns is REST API for campaigns
func AdminCampaigns(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	protected := protectedWithBasicAuth(w, r)
	if !protected {
		return
	}

	switch r.Method {
	case http.MethodPost:
		createCampaign(w, r)
	case http.MethodPut:
		updateCampaign(w, r)
	case http.MethodDelete:
		deleteCampaign(w, r)
	case http.MethodGet:
		listCampaigns(w, r)
	default:
		listCampaigns(w, r)
	}
}

func protectedWithBasicAuth(w http.ResponseWriter, r *http.Request) bool {
	// Protect endpoint with basic auth
	authUser, authPass, _ := r.BasicAuth()
	if !(strings.EqualFold(os.Getenv("AUTH_USER"), authUser) &&
		strings.EqualFold(os.Getenv("AUTH_PASSWORD"), authPass)) {
		http.Error(w, `{"status":"error", "msg": "Missing or wrong credentials."}`, http.StatusForbidden)
		return false
	}

	return true
}


func createCampaign(w http.ResponseWriter, r *http.Request) {
	repository := Repository{appengine.NewContext(r)}

	jsonFromBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorToJSON(w, err)
		return
	}

	var campaign Campaign
	err = json.Unmarshal(jsonFromBody, &campaign)
	if err != nil {
		ErrorToJSON(w, err)
		return
	}

	_, err = repository.SaveCampaign(&campaign, 0)
	if err != nil {
		ErrorToJSON(w, err)
		return
	}

	json, _ := json.Marshal(campaign)
	w.Write(json)
	return
}

func updateCampaign(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[(len("/admin/campaigns") + 1):])
	if err != nil {
		ErrorToJSON(w, err)
		return
	}

	repository := Repository{appengine.NewContext(r)}

	jsonFromBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorToJSON(w, err)
		return
	}

	campaign, err := repository.GetCampaign(int64(id))
	if err != nil {
		ErrorToJSON(w, err)
		return
	}

	err = json.Unmarshal(jsonFromBody, &campaign)
	if err != nil {
		ErrorToJSON(w, err)
		return
	}

	_, err = repository.SaveCampaign(&campaign, int64(id))
	if err != nil {
		ErrorToJSON(w, err)
		return
	}

	json, _ := json.Marshal(campaign)
	w.Write(json)
	return
}

func listCampaigns(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/admin/campaigns/") {
		getCampaign(w, r)
		return
	}

	repository := Repository{appengine.NewContext(r)}

	campaigns, err := repository.ListCampaigns(r.URL.Query())
	if err != nil {
		ErrorToJSON(w, err)
		return
	}

	json, _ := json.Marshal(campaigns)
	w.Write(json)
	return
}

func getCampaign(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/admin/campaigns/"):])
	if err != nil {
		ErrorToJSON(w, err)
		return
	}

	context := appengine.NewContext(r)
	repository := Repository{context}
	campaign, err := repository.GetCampaign(int64(id))
	if err != nil {
		ErrorToJSON(w, err)
		return
	}

	json, _ := json.Marshal(campaign)
	w.Write(json)
	return
}

func deleteCampaign(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[(len("/admin/campaigns") + 1):])
	if err != nil {
		ErrorToJSON(w, err)
		return
	}

	repository := Repository{appengine.NewContext(r)}
	err = repository.DeleteCampaign(int64(id))
	if err != nil {
		ErrorToJSON(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
