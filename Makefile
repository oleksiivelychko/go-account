go-test:
	$(info Run tests sequentially with limit of the maximum number of test running in parallel)
	$(warning main:19 env.LoadEnv - uncomment for local testing)
	$(warning oleksiivelychko/go-jwt-issuer application must be running before)
	go clean -testcache && go test ./*/ -p 1

heroku-bash:
	heroku run bash -a oleksiivelychkogoaccount

heroku-logs:
	heroku logs -n 200 -a oleksiivelychkogoaccount --tail
