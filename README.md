# go-account

### Microservice provides API to manage users accounts including ACL, enhanced by third-party JWT microservice.

📌 There are available environment variables with default values:
```
HOST=localhost
PORT=8081
DB_LOG=disable
DB_HOST=localhost
DB_PORT=5432
DB_NAME=account
DB_USER=admin
DB_PASS=secret
DB_DRIVER=postgres
DB_SSL=disable
DB_TZ=UTC
TEST_DB_HOST=localhost
TEST_DB_PORT=5433
TEST_DB_NAME=account-test
TEST_DB_USER=test
TEST_DB_PASS=secret
APP_JWT_URL=http://localhost:8080
```

💡 Watch logs of app for single pod:
```
kubectl exec goaccount-pod-0 -n gons -- tail -f /app/logs/gorm_$(date +%d-%m-%Y).log -n 100
```

💡 <a href="https://github.com/oleksiivelychko/go-jwt-issuer">JWT issuer app</a> must be running before.

![How it works](social_preview.png)
