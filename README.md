[![Build Status](https://travis-ci.org/I-Dont-Remember/deals-api.svg?branch=master)](https://travis-ci.org/I-Dont-Remember/deals-api)
[![Coverage Status](https://coveralls.io/repos/github/I-Dont-Remember/deals-api/badge.svg?branch=master)](https://coveralls.io/github/I-Dont-Remember/deals-api?branch=master)
# Deals API

API for bar locations and deals in Madison, WI  

Originally part of [Gimme Deals](https://github.com/I-Dont-Remember/GimmeDeals) but have since
switched to separate repos.

Deprecated directory will last until we feel it isn't useful as a reference point.

### Deployment / Serverless
This API uses the Serverless framework to make provisioning new Lambda/API Gateway setups or updating them a breeze.  We could loop in creating the DynamoDB under resources, but for persistent things like data it seems like its' better to keep it separate because otherwise CloudFormation(which is the magic of Serverless) gets confused because it already exists and fails.  

### Upload
First make sure the tables exist (`/scripts/create-dynamo-tables.sh`)  
Use a compiled binary or `go run` with `tools/upload/upload.go` and the input files directory. Use the `toml.example` file for how to structure them.  

### Versioning
Once we have gotten past the initial stage, we will be using SemVer as it is fairly common and easy to grasp. Specifics of how it will be gated will be documented once completed.

### Endpoints (and the connected Lambda functions)

#### /deals (GET)
_/functions/deals/all_  
Desc: Acquire deals  
Query Options
  - Day=[M,Tu,W,Th,F,Sa,Su]
  - Location=\<ID string\>
  - (Done by client for now) Time

#### /locations (GET)
_/functions/locations/all_  
Desc: Acquire locations

#### /locations/[name] (GET)
_/functions/locations/getID_  
Desc: Get location ID from name

### CI
To setup with Travis CI, need AWS credentials as env variables; as well as a GitHub Personal Access Token since netrc uses HTTPS. 