# o7tracker

o7tracker is high-performance `click tracker` designed for [Google App Engine][gae] written in [Go].

[![Build Status][travis-ci-badge]][travis-ci]
[![Go Report Card][goreportcard-badge]][goreportcard]

## Development & Setup

Make sure you have [Google App Engine SDK for Go][gae-sdk-go] installed. Then you can run:

```bash
goapp serve
```

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
GET /admin/campains/:id
GET /admin/campains?platforms=Android,WindowsPhone
PUT /admin/campains/:id
DELETE /admin/campains/:id
```


## Assumptions

- This project assumes that Go is the fastest citizen of [Google App Engine][gae] for this particular task.
- This project uses that layer caching with sequence `global variable > memcache > datastore`.

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
