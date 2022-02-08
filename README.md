# go-account

### Microservice provides API to manage users accounts including ACL through third-party JWT microservice.

💡 Deployed on <a href="https://oleksiivelychkogoaccount.herokuapp.com">Heroku</a>

Run tests sequentially with limit of the maximum number of test running in parallel:
```
go clean -testcache && go test ./*/ -p 1
```

There are next environment variables available with default values:
```
PORT=8081
DB_LOG=enable
DB_HOST=localhost
DB_PORT=5432
DB_NAME=go-postgres
DB_USER=gopher
DB_PASS=secret
DB_DRIVER=postgres
DB_SSL=require
DB_TZ=UTC
TEST_DB_HOST=localhost
TEST_DB_PORT=5433
TEST_DB_NAME=go-postgres-test
TEST_DB_USER=gopher
TEST_DB_PASS=secret
APP_JWT_URL=http://0.0.0.0:8080
```

**P.S.** <a href="https://github.com/oleksiivelychko/go-jwt-issuer">JWT issuer</a> must be running before.

Use `APP_JWT_URL=https://oleksiivelychkogojwtissuer.herokuapp.com` as remote backend service.

![how it works](.dock/readme.png)
