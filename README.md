[![Shortener test](https://github.com/amanakin/shortener/actions/workflows/test.yaml/badge.svg?branch=master)](https://github.com/amanakin/shortener/actions/workflows/test.yaml)

## Yet Another URL Shortener

### Description
This is a simple URL shortener service.\
It provides HTTP and gRPC API.\
It uses PostgreSQL/in-memory as a storage.\

### Build/Deploy
Most convenient way to deploy is to use Docker.

For running `shortener` server just execute:
```shell
docker compose up -d
```
See .env.example for environment variables.\
(Do not use this example for deployment)

By default, service consumes config with -c flag.\
See [etc/shorter.yaml](etc/shortener.yaml):
```yaml
shortener:
  alphabet: "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
  short_len: 10
  default_scheme: "https"
  allowed_schemes:
    - "http"
    - "https"
http:
  enabled: true
  host: 0.0.0.0
  port: 8080
  rate_limit: 100
grpc:
  enabled: true
  host: 0.0.0.0
  port: 8081
postgres:
  enabled: true # false to use in-memory storage
  host: postgresdb
  port: 5432
  user: dev
  password: dev_password
  dbname: shortener

```

To build server:
```shell
make all
```

To update mock:
```shell
make mockgen
```

To update protobuf and grpc:
```shell
make protogen
```

To run unit-tests:
```shell
make test
```

### API

The service provides the ability to shorten links and get a original link for shortened one.\
Actually it returns unique token, which should be use by redirect microservice.

For using see HTTP [swagger](api/http/shortener.yaml) and [protobuf](api/grpc/shortener.proto) API.

To play with HTTP I recommend [Postman](https://www.postman.com/)\
To play with gRPC you can try:
```shell
go run cmd/grpc-client/main.go
```
Example:
```
shorten google.com
original: https://google.com
shortened: mCMdDvvigK
created: false
resolve mCMdDvvigK
resolved: https://google.com
shorten ftp://impossible.com
error: rpc error: code = Unknown desc = shorten: invalid URL
shorten http://my.com
original: http://my.com
shortened: kOztwSe9OX
created: true
```


