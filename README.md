# go-account

### Microservice provides API to manage users accounts including ACL through third-party JWT microservice.

⚙️ Deployed on <a href="https://oleksiivelychkogoaccount.herokuapp.com">Heroku</a>

There are available environment variables with default values:
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
DATABASE_URL=postgres://gopher:secret@host:port/go-postgres
TEST_DB_HOST=localhost
TEST_DB_PORT=5433
TEST_DB_NAME=go-postgres-test
TEST_DB_USER=gopher
TEST_DB_PASS=secret
APP_JWT_URL=http://0.0.0.0:8080
```

💡 <a href="https://github.com/oleksiivelychko/go-jwt-issuer">JWT issuer app</a> must be running before.

![how it works](.http-requests/readme.png)
