[+] Deletion Endpoint
[+] Deletion Endpoint Tests
[ ] KGS
    [+] Caching Mechanism
        [+] Saves Newly Generated Alias For 24 Hours into Cache
        [+] On Read First Checks Cache
        [ ] Unit Tests
    [+] MonoRepo
        [+] Separate Main Service
    [+] Add KGS to MonoRepo
[ ] Analytics Storage

Start Postgres Container
```shell
docker run --name alias-gen-postgres -e POSTGRES_USER=alias-gen -e POSTGRES_PASSWORD=alias-gen -e POSTGRES_DB=url-aliases -d -p 5432:5432 postgres
```