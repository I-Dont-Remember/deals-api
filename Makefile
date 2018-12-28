.PHONY: deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/not-implemented functions/not-implemented.go

test-all:
	(cd functions && go test ./...)

deploy: build
	npm run sls -- deploy --verbose
