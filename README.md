# o7tracker

o7tracker is high-performance click tracker designed for [Google App Engine][gae] written in [Go].

[![Build Status][travis-ci-badge]][travis-ci]

## Development & Setup

Make sure you have [Google App Engine SDK for Go][gae-sdk-go] installed. Then you can run:

```bash
AUTH_USER=admin AUTH_PASSWORD=admin goapp serve -clear_datastore
```

- `AUTH_USER` and `AUTH_PASSWORD` are settings for HTTP basic authentication.
- For producton environment please [see app.yaml](./app.yaml).
- `-clear_datastore` flag make sure that local Datastore is clean.

And deploy with

```bash
appcfg.py -A <app_name> -V v1 update ./
```

Tracker should be up and running on 

```bash
curl <app_name>.appspot.com
```

## APIs

### Tracker API

Tracker exposes one endpoint that requires presence of `campaign_id` and `platform`. 

```
GET /track/:id?platform=Android
```

### Admin API

[Admin campaigns API](./admin_campaigns.go) requires [HTTP basic authentication][basic-auth] on all following endpoints:

```
POST /admin/campains
```

With following JSON payload:

```json
{
    "redirect_url":"https://github.com/otobrglez/o7tracker",
    "platforms":[
        "WindowsPhone", "Android", "IPhone"
    ]
}
```

Result:

```json
{
  "id": 5928566696968192,
  "redirect_url": "https://github.com/otobrglez/o7tracker",
  "created_at": "2016-04-10T21:47:44.519184234Z",
  "updated_at": "2016-04-10T21:47:44.519184189Z",
  "platforms": [
    "WindowsPhone",
    "Android",
    "IPhone"
  ],
  "click_count": 0,
  "android_click_count": 0,
  "iphone_click_count": 0,
  "windowsphone_click_count": 0
}
```

Other endpoints:
```
GET /admin/campains/:id
GET /admin/campains?platforms=Android,WindowsPhone
PUT /admin/campains/:id
DELETE /admin/campains/:id
```


Help yourself with following Postman examples:

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/be2b92e7ffc18a31ae38)

## Assumptions

- This project assumes that Go is the fastest citizen of [Google App Engine][gae] and fits perfectly for this particular task.
- This project uses layered caching with sequence `memcache > datastore`.
- Further optimisation could also be to parse `User-Agent` from browser
instead of relaying on `?platform=`
- Further optimisation could be to stream clicks into `InfluxDB` to get
real-time overview of clicks.
- Further optimisation in caching could also introduce global
slice variable. However that would require some additional work to
support cases like deletion of campaigns.


## Testing

- Code comes with test suite for main components in [repository](./repository_test.go)
- and [integration test](./integration_test.sh) written in Bash.


## Author

- [Oto Brglez](https://github.com/otobrglez)


[go]: https://golang.org/
[gae]: https://cloud.google.com/appengine/ 
[gae-sdk-go]: https://cloud.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go
[travis-ci]: https://travis-ci.org/otobrglez/o7tracker
[travis-ci-badge]: https://travis-ci.org/otobrglez/o7tracker.svg?branch=master
[goreportcard-badge]: https://goreportcard.com/badge/otobrglez/o7tracker
[goreportcard]: https://goreportcard.com/report/otobrglez/o7tracker
[basic-auth]: https://en.wikipedia.org/wiki/Basic_access_authentication
