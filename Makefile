go-test:
	$(info Run tests sequentially with limit of the maximum number of test running in parallel.)
	go clean -testcache && go test ./*/ -p 1

docker-build:
	[[ -z "$(docker images -q local/goaccount)" ]] || docker image rm local/goaccount
	docker build --tag local/goaccount .

docker-push: docker-build
	$(warning instead of `local` prefix use dockerhub account name and change/remove `imagePullPolicy`)
	docker buildx build --platform linux/amd64 --tag local/goaccount .
	docker push local/goaccount
