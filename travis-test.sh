#! /bin/bash

cd ./scripts
bash test-scripts.sh
cd ..

cd ./tools/upload/ && go test -v \
&& cd ../../functions/ && go test -v ./...