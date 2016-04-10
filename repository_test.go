package o7tracker

import (
	"appengine/aetest"
	"appengine/datastore"
	"fmt"
	"testing"
)

type BankAccount struct {
	Amount int
}

func TestAddCampaign(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	key := datastore.NewKey(c, "BankAccount", "", 1, nil)
	if _, err := datastore.Put(c, key, &BankAccount{100}); err != nil {
		t.Fatal(err)
	}

	fmt.Println(key)
}

func TestTrack(t *testing.T) {

}
