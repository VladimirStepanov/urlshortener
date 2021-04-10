# Description

Simple url-shortener written in Golang with Redis database

## Build

```bash
docker-compose build
```

## Test

Assume that Redis already running on port 6379

```bash
make test
```

## Run

```bash
docker-compose up
```

# Endpoints

## Get encoded URL info

`GET /info/{encoded_url}`

```bash
curl -L -X GET 'http://localhost:8080/WuYbydedVqi'
```

### Response

```json
{
    "id":"WuYbydedVqi",
    "url":"https://www.alexedwards.net/blog/working-with-redis",
    "visits":2,
    "expire":"4.10.2022 17:18:0",
    "once":false
}
```

## Encode URL

`POST /encode`

Params (json):
* url [string]
* expire - UTC date in format d.m.y h:m:s [string]
* once - allows only one redirect  [boolean]

```bash
curl -L -X POST 'localhost:8080/encode' -H 'Content-Type: application/json' --data-raw '{
    "url": "https://www.alexedwards.net/blog/working-with-redis",
    "expire": "4.10.2022 17:18:00",
    "once": false
}'
```

### Response

```json
{
    "status":"success",
    "url":"http://localhost:8080/YbnuLt4L5Eu"
}
```

## Redirect to original URL

`GET /{encoded_url}`

```bash
curl -L -X GET http://localhost:8080/OTv0FdGU8Ng
```

## Delete encoded URL

`DELETE /{encoded_url}`

```bash
curl -L -X DELETE http://localhost:8080/OTv0FdGU8Ng
```
