# GoApp

API endpoint on Go server for processing device events and saving it to Clickhouse.
Events are sent to clickhouse in batches (default=1000) or if no income data for 5 seconds

## Requirements
- docker
- docker-compose


## Instructions

- Clone the repository and move to the project directory:
  ```bash
    git clone https://github.com/vbelorus/go-app.git
    cd go-app
  ```
- Build images
      ```bash
        docker-compose build
      ```
- Review env variables in .env
    ```
    HTTP_SERVER_PORT=8080
    CLICKHOUSE_HOST=clickhouse-server
    CLICKHOUSE_PORT=9000
    CLICKHOUSE_DATABASE=app
    CLICKHOUSE_USERNAME=default
    CLICKHOUSE_PASSWORD=
    DEVICE_EVENT_BATCH_SIZE=1000
    ```

- Run containers
   ```bash
     docker-compose up -d
   ```
Go server will up on :8080 by default

## Usage Examples
-  Add Event Example
  ```bash
    curl -i -X POST 127.0.0.1:8080/events -d \
      '{
        "client_time":"2020-12-01 23:59:00",
        "device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E",
        "device_os":"iOS 13.5.1",
        "session":"ybuRi8mAUypxjbxQ",
        "sequence":1,
        "event":"app_start",
        "param_int":0,
        "param_str":"some text"
      }'
  ```
Response should be:
    ```HTTP/1.1 200 OK```
