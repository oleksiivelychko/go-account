go-test:
	$(info Run tests sequentially with limit of the maximum number of test running in parallel)
	go clean -testcache && go test ./*/ -p 1
