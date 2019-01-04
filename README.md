# Deals API

**Don't forget [YAGNI](https://www.martinfowler.com/bliki/Yagni.html)**

A (hopefully) self contained and well documented repo for Madtown Deals/ Deals On Tap/ or whatever it's called now. The various functions are stored under `functions` by the description of their operation.  They are all implemented as standalone handlers that take a request & a DB reference, so it would be trivial to set them all up as functions from a single function (in the case we are seeing issues with cold starts, etc), though we would then have to handle our own routing & other stuff that API Gateway currently handles.  The functions handling the logic for each request, not just the Lambda stuff, are in multiple files at the top of the model directory such as `locations/createlocation.go` under a single package, such as `location`.  That way it's easy to import them into Lambda handlers, tests, local server handling, and maybe more.

A quick and dirty form of auth is done with `helpers.AuthMiddleware`.  It checks if the `X-Dot-Auth` header matches the value in the environment variable `API_AUTH`.  This can be thrown away once it becomes time to implement an actual auth provider.


## Formatting & Linting
Currently developed in VS Code with the Go by Microsoft extension.  This uses `gofmt`, `golint`, and several other tools to help with development and save time formatting your code and checking for issues.

## Development

To get setup you must have Node & Go setup on your machine.  For Node, it is recommended to use `nvm` to easily install different versions without a lot of effort.  Until a concerted effort is made to make this simpler, Go should be installed through your default package manager if possible, otherwise follow the regular instructions to install `v1.10.4`. Once you have both, run `npm i` in the repo to install the necessary tools into your directory (just Serverless Framework for now, installing in repo to help keep everything necessary for dev in one place).  To install Go dependencies, try `go get -t ./...`, but no guarantees that will work correctly.  Eventually that will get fixed by trying to use Go modules or dep tool or etc.

**TODO:** handle versioning of Go dependencies & the language itself so multiple developers aren't out of sync.

To be able to deploy the application, you must have AWS credentials setup with the `aws-cli` tool.

An environment variable is used to differentiate between the different options.`API_ENV` can select from either `local`: for using [localstack](https://github.com/localstack/localstack), `prod`: for deployments on AWS, or `test`: for running tests with the Mock DB. 

Since there is no serverless-offline plugin for Go, to check functionality we can use a locally running DynamoDB and come fairly close to how AWS runs.  By keeping most of the logic inside packages, it is easy to import and pass requests to the logic functions. `go run local.go` connects while serving the API on `localhost`.  It uses the Echo web framework and a function that massages the request/response objects from Echo to match the inputs and outputs of the Lambda handler functions.  While convenient for testing, it does need a vigilant watch to make sure it matches the expected url paths & other things that API Gateway/Serverless handle.

At some point this can be extended to write fairly decent integration tests, but for now it's effective for manual testing.  One of the benefits of it being a little more difficult to test locally than some random script is that it pushes TDD (Test Driven Development), because running `go test` is significantly faster than the using `local` or `prod` environments.

## Testing

To test a package or function, cd to it's directory and run `go test`.
To run the entire test suite, `cd functions && go test ./...`.  This will show many directories without test files, as expected.
