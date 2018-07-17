# HTTP service monitor

Monitor your APIs, services, web sites and pretty much everything which works via http and json.
Set up tests, customize execution schedule, collect statistics and report results through various channels.

[![Build Status](https://travis-ci.org/viktorminko/monitor.svg?branch=master)](https://travis-ci.org/viktorminko/monitor)

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

```
./monitor -workdir your_config_directory
```
 
 For other usage and deployment options see https://github.com/viktorminko/monitor_tools
