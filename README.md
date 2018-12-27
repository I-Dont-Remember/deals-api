# Deals API

A (hopefully) self contained and well documented repo for Madtown Deals/ Deals On Tap/ or whatever it's called now. The various functions are stored under `functions` by the description of their operation.  They are all implemented as standalone handlers that take a request & a DB reference, so it would be trivial to set them all up as functions from a single function (in the case we are seeing issues with cold starts, etc), though we would then have to handle our own routing & other stuff that API Gateway currently handles.



## Development

To get setup you must have Node & Go setup on your machine.  Once you do, run `npm i` to install the necessary tools into your directory (just Serverless Framework for now, installing in repo to help keep eerything necessary for dev in one place).

To be able to deploy the application, you must have AWS credentials setup with the `aws-cli` tool.

Since there is no serverless-offline plugin for Golang, to check functionality we can use a test suite connected to a locally running DB and come fairly close to how AWS runs, though we have to build the request in the test rather than making actual HTTP calls.  An environment variable is used to differentiate between the different options.`API_ENV` can select from either `local`: for using [localstack](https://github.com/localstack/localstack), `prod`: for deployments on AWS, or `test`: for running tests with the Mock DB.

To test a function, cd to it's directory and run `API_ENV=<test or local> go test`.  To run it with `local` you must have started Localstack and run the create tables script.