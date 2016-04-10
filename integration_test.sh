#!/usr/bin/env bash
set -e

BASE_URI="http://localhost:8080"
CPARAMS="-D -"
HTTP_USER=admin
HTTP_PASSWORD=admin

echo -e "\n=====> Basic CRUD campaigns test"

echo "GET /"
curl -s $CPARAMS $BASE_URI/ \
  | grep -q outfit7 || exit 1

echo "GET /track"
curl -s -L --max-redirs 1 $CPARAMS $BASE_URI/track \
  | grep -q outfit7 || exit 1

echo "GET /track/"
curl -s $CPARAMS $BASE_URI/track/ \
  | grep -q outfit7 || exit 1

echo "POST / (no access)"
curl -s \
  -X POST \
  -H "Accept: application/json" \
  -d "{}" \
  "$BASE_URI/admin/campaigns" \
  | grep -q -i "Missing" || exit 1

echo -e "\n=====> Admin / campaigns"

RANDOM_STRING=$(hexdump -n 16 -v -e '/1 "%02X"' -e '/16 "\n"' /dev/urandom)
UPDATE_URL="http://github.com/otobrglez/o7tracker?pom=$RANDOM_STRING"

echo "POST /admin/campaigns"
CAMPAIGN_ID=$(cat <<JSON |
{
    "redirect_url":"$UPDATE_URL",
    "platforms":[
        "WindowsPhone", "Android", "IPhone"
    ]
}
JSON
curl -s \
  -u $HTTP_USER:$HTTP_PASSWORD \
  -X POST \
  -H "Accept: application/json" \
  -d @- "$BASE_URI/admin/campaigns" | jq -r ".id")

echo "GET /admin/campaigns/$CAMPAIGN_ID"
curl -s \
  -u $HTTP_USER:$HTTP_PASSWORD \
  -H "Accept: application/json" \
  $BASE_URI/admin/campaigns/$CAMPAIGN_ID \
  | jq -r ".id" | grep -q $CAMPAIGN_ID || exit 1


echo "PUT /admin/campaigns/$CAMPAIGN_ID"

RANDOM_STRING=$(hexdump -n 16 -v -e '/1 "%02X"' -e '/16 "\n"' /dev/urandom)
UPDATE_URL="http://github.com/otobrglez/o7tracker?pom=$RANDOM_STRING"

cat <<JSON |
{
    "redirect_url":"$UPDATE_URL",
    "platforms":["WindowsPhone","Android"]
}
JSON
curl -s \
  -u $HTTP_USER:$HTTP_PASSWORD \
  -X PUT \
  -H "Accept: application/json" \
  -d @- "$BASE_URI/admin/campaigns/$CAMPAIGN_ID" | \
  jq -r ".redirect_url" \
  | grep -q $RANDOM_STRING || exit 1

sleep 0.5

echo "GET /admin/campaigns?platform=WindowsPhone,IPhone"
curl -s \
  -u $HTTP_USER:$HTTP_PASSWORD \
  -H "Accept: application/json" \
  $BASE_URI/admin/campaigns?platform=WindowsPhone,IPhone \
  | grep -q $RANDOM_STRING || exit 1

echo -e "\n=====> Redirection Test"

curl -s -L --max-redirs 1 $CPARAMS \
  $BASE_URI/track/$CAMPAIGN_ID?platform=IPhone | grep -q "github" || exit 1

curl -s -L --max-redirs 1 $CPARAMS \
  $BASE_URI/track/$CAMPAIGN_ID?platform=IPhone | grep -q "github" || exit 1

curl -s -L --max-redirs 1 $CPARAMS \
  $BASE_URI/track/$CAMPAIGN_ID?platform=IPhone | grep -q "github" || exit 1

curl -s -L --max-redirs 1 $CPARAMS \
  $BASE_URI/track/$CAMPAIGN_ID?platform=WindowsPhone | grep -q "github" || exit 1

sleep 1

curl -s \
  -u $HTTP_USER:$HTTP_PASSWORD \
  -H "Accept: application/json" \
  $BASE_URI/admin/campaigns/$CAMPAIGN_ID \
  | jq ".click_count" \
  | grep -q 4 || exit 1

echo "Done with suite."
exit 0

echo "DELETE /admin/campaigns/$CAMPAIGN_ID"
curl -s \
  -u $HTTP_USER:$HTTP_PASSWORD \
  -X DELETE \
  -H "Accept: application/json" \
  $BASE_URI/admin/campaigns/$CAMPAIGN_ID \

sleep 0.5

echo "GET /admin/campaigns"
curl -s \
  -u $HTTP_USER:$HTTP_PASSWORD \
  -H "Accept: application/json" \
  $BASE_URI/admin/campaigns \
  | grep -q -v $CAMPAIGN_ID || exit 1

echo "Done wit cleanup."
