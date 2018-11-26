# HTTP service monitor

Monitor your APIs, services, web sites and pretty much everything which works via http and json.
Set up tests, customize execution schedule, collect statistics and report results through various channels.

[![Build Status](https://travis-ci.org/viktorminko/monitor.svg?branch=master)](https://travis-ci.org/viktorminko/monitor)
[![Coverage Status](https://coveralls.io/repos/github/viktorminko/monitor/badge.svg?branch=master&service=github)](https://coveralls.io/github/viktorminko/monitor)

## Features

### Flexible scheduling
Customize executions periods for every test, whole test suite and statistics reporting

### Authorization support
Application can perform access tokens requests to specific authorization URL and add tokens to test requests

### Statistics collection
Monitor collects different types of tests statistics: number of success/failed requests, average response time etc.

### Various notification channels
Test errors and statistic reports can be send by email or telegram, so you can get immediate notification if something goes wrong
 
## Usage

With compileв binary

```
./monitor -workdir your_config_directory
```
 
 With docker image
```
docker run -v $(pwd)/config:/app/config viktorminko/monitor
``` 
