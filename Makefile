build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/get-deals functions/deals/get/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/post-deals functions/deals/post/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get-locations functions/locations/get/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/post-locations functions/locations/post/main.go

.PHONY: deploy
deploy: build
	sls deploy --verbose
