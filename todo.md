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
- [x] Analytics Storage with UI
    - [x] Base Event Model
    - [x] Add ClickHouse
    - [x] Send Event To ClickHouse
    - [x] Extend Event
    - [x] Add Metabase UI
- [x] Refactor
- [x] Dockerize
- [ ] Commands in Makefile

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
