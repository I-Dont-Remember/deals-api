.PHONY: deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/not-implemented functions/not-implemented/main.go

deploy: build
	npm run sls -- deploy --verbose
