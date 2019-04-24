.PHONY: deploy build

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/not-implemented functions/not-implemented.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/analytics functions/analytics/post/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/email functions/email/post/main.go
	for dir in functions/deals/*/;  do \
		echo "$$dir"; \
		name="$$(basename $$dir)"; \
		echo "Building $$name..."; \
		env GOOS=linux go build -ldflags="-s -w" -o bin/"$$name" "$$dir"/main.go; \
		done
	for dir in functions/locations/*/;  do \
		echo "$$dir"; \
		name="$$(basename $$dir)"; \
		echo "Building $$name..."; \
		env GOOS=linux go build -ldflags="-s -w" -o bin/"$$name" "$$dir"/main.go; \
		done
	for dir in functions/campuses/*/;  do \
		echo "$$dir"; \
		name="$$(basename $$dir)"; \
		echo "Building $$name..."; \
		env GOOS=linux go build -ldflags="-s -w" -o bin/"$$name" "$$dir"/main.go; \
		done

test-all:
	(cd functions && go test ./...)

dev-deploy: build
	if test -z "$$API_AUTH"; then { echo "API_AUTH not set"; exit 1; } else (npm run sls -- deploy --stage dev --verbose) fi


prod-deploy: build
	if test -z "$$API_AUTH"; then { echo "API_AUTH not set"; exit 1; } else (npm run sls -- deploy --stage prod --verbose) fi
