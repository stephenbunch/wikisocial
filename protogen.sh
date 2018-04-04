#!/bin/bash

rm -rf server/_proto
rm -rf web/lib/_proto
mkdir -p server/_proto
mkdir -p web/lib/_proto

protoc \
--plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go \
--plugin=protoc-gen-ts=web/node_modules/.bin/protoc-gen-ts \
--go_out=plugins=grpc:server/_proto \
--js_out=import_style=commonjs,binary:web/lib/_proto \
--ts_out=service=true:web/lib/_proto \
--proto_path=./proto \
wikitribe.proto
