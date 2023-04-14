docker-build:
	[[ -z "$(docker images -q local/goaccount)" ]] || docker image rm local/goaccount
	docker build --no-cache --tag local/goaccount .

docker-push: docker-build
	$(warning instead of `local` prefix use dockerhub account name and change/remove `imagePullPolicy`)
	docker buildx build --platform linux/amd64 --tag local/goaccount .
	docker push local/goaccount

docker-network:
	docker network inspect go-network >/dev/null 2>&1 || docker network create --driver bridge go-network

run:
	HOST=localhost \
	PORT=30000 \
	DB_HOST=localhost \
	DB_PORT=5432 \
	DB_NAME=account \
	DB_USERNAME=admin \
	DB_PASSWORD=secret \
	DB_TIMEZONE=UTC \
	DB_SSL_MODE=disable \
	DB_LOG_PATH=.log \
	APP_JWT_URL=http://localhost:8080 \
	go run main.go

test:
	$(info Run tests sequentially with limit of the maximum number of test running in parallel.)
	go clean -testcache && go test ./*/ -p 1

log:
	kubectl exec goaccount-pod-0 -n gons -- tail -f /app/logs/gorm_$(date +%d-%m-%Y).log -n 100
