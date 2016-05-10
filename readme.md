# Tarzan
Awwwwoooo, It's my personal project to scraping themef****t.net via tor. We're watching all items they have, the old and the new. We collects items' data such as today sales, tags, etc.

## Requirements
 - MongoDB ([www.mongodb.org](https://www.mongodb.org))
 - Redis ([redis.io](http://redis.io/))
 - Tor (Run tor with config .torrc, geoip and geoip6 from ```dist``` directory)

## Build
Run ```build.sh``` and you'll got 3 binaries:

- ```api``` will run API server on port 8080, with endpoints url ```/api/v1```
- ```worker``` will always run and retrieve task from redis via ```RPUSH```
- ```scheduler``` will crawling category/subcategory page via API, where we can manage it via web API, it's like cronjob