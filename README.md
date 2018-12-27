# Deals API

A (hopefully) self contained and well documented repo for Madtown Deals/ Deals On Tap/ or whatever it's called now.


## Development

To get setup you must have Node & Go setup on your machine.  Once you do, run `npm i` to install the necessary tools into your directory (just Serverless Framework for now, installing in repo to help keep eerything necessary for dev in one place).

To be able to deploy the application, you must have AWS credentials setup with the `aws-cli` tool.

Since there is no serverless-offline plugin for Golang, to check functionality we can use a test suite connected to a locally running DB and come fairly close to how AWS runs, though we have to build the request in the test rather than making actual HTTP calls.  An environment variable is used to differentiate between the different options.`API_ENV` can select from either `local`: for using [localstack](https://github.com/localstack/localstack), `prod`: for deployments on AWS, or `test`: for running tests with the Mock DB.