# demogods
[![Build Status](https://travis-ci.org/frankgreco/demogods.svg?branch=master)](https://travis-ci.org/frankgreco/demogods) [![Docker Repository on Quay](https://quay.io/repository/frankgreco/demogods/status "Docker Repository on Quay")](https://quay.io/repository/frankgreco/demogods)

> a demo go application with prometheus instrumentation.

## Quick Start
```sh
$ mkdir -p $GOPATH/src/frankgreco
$ cd $GOPATH/src/frankgreco
$ git clone git@github.com:frankgreco/demogods.git
$ cd demogods
$ make
$ ./demogods &
$ curl http://localhost:8080
You live demo will succeed!
```
## Docker
```sh
docker run -d -p 8080:8080 -p 9000:9000 quay.io/frankgreco/demogods
$ curl http://localhost:8080
You live demo will succeed!
```

## Prometheus
Multiple metrics are exposed by this demo application. These metrics can be scraped by Prometheus at the following endpoint.
```
$ curl http://localhost:9000
# HELP demogods_demo_bad_count Count of all live demos that will fail.
# TYPE demogods_demo_bad_count counter
demogods_demo_bad_count 0
# HELP demogods_demo_good_count Count of all live demos that will be successful.
# TYPE demogods_demo_good_count counter
demogods_demo_good_count 1
...
```