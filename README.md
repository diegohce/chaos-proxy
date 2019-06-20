# chaos-proxy

[![Go Report Card](https://goreportcard.com/badge/github.com/diegohce/chaos-proxy)](https://goreportcard.com/report/github.com/diegohce/chaos-proxy)
[![GitHub release](https://img.shields.io/github/release/diegohce/chaos-proxy.svg)](https://github.com/diegohce/chaos-proxy/releases/)
[![Github all releases](https://img.shields.io/github/downloads/diegohce/chaos-proxy/total.svg)](https://github.com/diegohce/chaos-proxy/releases/)
[![GPLv3 license](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://github.com/diegohce/chaos-proxy/blob/master/LICENSE)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/diegohce/chaos-proxy/graphs/commit-activity)
[![HitCount](http://hits.dwyl.io/diegohce/chaos-proxy.svg)](http://hits.dwyl.io/diegohce/chaos-proxy)
[![Generic badge](https://img.shields.io/badge/deb%20package-no-red.svg)](https://github.com/diegohce/chaos-proxy/releases/)

## What is it?
Chaos-proxy acts almost as a normal proxy. It causes at random, normal http request error scenarios. 500's errors, timeout delays and connections drop.

## But why?
I use it to test microservices resilience. If it works fine with chaos-proxy in the middle, it will certainly work without it.

## Bind address

Default bind address and port: `0.0.0.0:6667`. It can be changed setting `CHAOSPROXY_BINDADDR` environment variable.


## Config file

`chaos-proxy.json` can be placed into project directory or, preferably, in `/etc/chaos-proxy`

* max_timeout: Sets the top boundary for random milliseconds to timeout the request.
* default_host: Where every request that does not match one of `paths` will be routed.
* paths: Where to route specific requests. If path ends with `/` it's interpreted as "begins with".

```json
{
	"max_timeout": 5000,
	"default_host": {
		"host": "http://localhost:6666"
	},
	"paths": {
		"/badservice/status/400": {
			"host": "http://localhost:6666"
		},
		"/badservice/status/403": {
			"host": "http://localhost:6666"
		},
		"/badservice/status/404": {
			"host": "http://localhost:6666"
		},
		"/data/2.5/": {
			"host": "http://api.openweathermap.org:80"
		}
	}
}
```

# Status

Chaos-proxy is still in a very early stage. The "random error" generator is pretty lousy and there's code that can be improved for sure. But, as is, it works as intended.


