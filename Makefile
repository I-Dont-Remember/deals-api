build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/get-deals functions/deals/get/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/post-deals functions/deals/post/main.go

.PHONY: deploy
deploy: build
	sls deploy --verbose
