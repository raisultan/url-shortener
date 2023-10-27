## Backlog
- [x] Deletion Endpoint
- [x] Deletion Endpoint Tests
- [x] KGS
    - [x] Caching Mechanism
        - [x] Saves Newly Generated Alias For 24 Hours into Cache
        - [x] On Read First Checks Cache
        - [ ] Unit Tests (🤷‍♀️)
    - [x] MonoRepo
        - [x] Separate Main Service
    - [x] Add KGS to MonoRepo
    - [x] Add KGS call to Main Service
- [ ] Analytics Storage with UI
    - [x] Base Event Model
    - [x] Add ClickHouse
    - [x] Send Event To ClickHouse
    - [x] Extend Event
    - [x] Add Metabase UI
- [ ] Dockerize

## Analytics
### Usage Analytics
- Total Clicks: Track the total number of clicks for each shortened URL
- Unique Clicks: Track the number of unique visitors who click on each shortened URL
- Click-Through Rate (CTR): Calculate the CTR by dividing the number of clicks by the number of times the shortened URL was shown
### User Analytics
- User Demographics: Collect demographic information about users, such as their location, device type, and browser
- Referrers: Track where the clicks are coming from (e.g., social media, search engines, direct traffic)
### Performance Analytics
- Latency: Measure the time it takes to redirect users from the shortened URL to the original URL
- Error Rates: Track the rate of errors, such as failed redirects or database errors
### Operational Analytics
- API Usage: Track the usage of your API, including the number of requests, errors, and the average response time
- Database Performance: Monitor the performance of your database, including query times and error rates

Storage: ClickHouse

UI: Metabase

## Useful Commands
Start Postgres Container
```shell
docker run --name alias-gen-postgres -e POSTGRES_USER=alias-gen -e POSTGRES_PASSWORD=alias-gen -e POSTGRES_DB=url-aliases -d -p 5432:5432 postgres
```

Docker Network For ClickHouse and Metabase
```shell
# create
docker network create url-shortener

# inspect
docker network inspect url-shortener
```

Start ClickHouse Container
```shell
docker run -d --name clickhouse-server \
    --ulimit nofile=262144:262144 \
    -p 9000:9000 -p 8123:8123 \
    -e CLICKHOUSE_DB=testing \
    -e CLICKHOUSE_SERVER__LISTEN_HOST='0.0.0.0' \
    --network url-shortener \
    yandex/clickhouse-server
```

Start Metabase Container
```shell
export METABASE_DOCKER_VERSION=v0.47.2
export METABASE_CLICKHOUSE_DRIVER_VERSION=1.2.2

mkdir -p mb/plugins && cd mb
curl -L -o plugins/ch.jar https://github.com/ClickHouse/metabase-clickhouse-driver/releases/download/$METABASE_CLICKHOUSE_DRIVER_VERSION/clickhouse.metabase-driver.jar
docker run -d -p 3000:3000 \
    --network url-shortener \
    --mount type=bind,source=$PWD/plugins/ch.jar,destination=/plugins/clickhouse.jar \
    metabase/metabase:$METABASE_DOCKER_VERSION
```
