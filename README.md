# o7tracker

o7tracker is `click tracker` designed for [Google App Engine][gae] written in [Go].


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

## API

### Tracker API

### Admin API

Admin API requires basic auth on all following endpoints:

```
POST /admin/campains
PUT /admin/campains/<id>
GET /admin/campains
DELETE /admin/campain/<id>
```

## Assumptions


## Author

- [Oto Brglez](https://github.com/otobrglez)


[go]: https://golang.org/
[gae]: https://cloud.google.com/appengine/ 
[gae-sdk-go]: https://cloud.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go
