### Local Development

In order to develop locally, you can use Makefile commands to set up dependencies and run services.

First, clone the repo and set up dependencies:

1. Clone the repository:
    ```bash
    git clone https://github.com/your-username/url-shortener.git
    ```
2. Navigate to the project directory:
    ```bash
    cd url-shortener
    ```
3. Run PostgreSQL container for `alias-gen` service:
    ```bash
    make run-postgres
    ```
4. Run Redis container:
    ```bash
    make run-postgres
    ```
5. Create Docker network for ClickHouse and Metabase:
    ```bash
    make create-docker-network
    ```
6. Run ClickHouse container:
    ```bash
    make run-clickhouse
    ```
7. Setup ClickHouse plugin for Metabase:
    ```bash
    make setup-metabase-plugins
    ```
8. Run Metabase container:
    ```bash
    make run-metabase
    ```

After setting up the dependencies you can create local configuration files `local.yaml` in `config`
directories of both services, they are pretty much the same as `production.yaml`. Also, instead of
using `mongo` as main storage you can use `SQLite`, to do that specify `active_storage: "sqlite"`
in `local.yaml` for `main` service and create storage directory in root:

```bash
mkdir storage
```

Finally, you can export `CONFIG_PATH` paths for both services and run them:

1. Export config path for `alias-gen` service and run it:
    ```bash
    export CONFIG_PATH=services/alias-gen/config/local.yaml && make run-alias-gen
    ```
2. Export config path for `main` service and run it:
    ```bash
    export CONFIG_PATH=services/main/config/local.yaml && make run-main
    ```
