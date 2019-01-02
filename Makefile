.PHONY: deploy build

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/not-implemented functions/not-implemented.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/create-location functions/locations/create-location/main.go

test-all:
	(cd functions && go test ./...)

deploy: build
	npm run sls -- deploy --verbose
