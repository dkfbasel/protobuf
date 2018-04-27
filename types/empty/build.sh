#!/bin/sh

protoc \
  --go_out=. \
  --proto_path=.. \
  empty.proto

mv ./github.com/dkfbasel/protobuf/types/empty/* .
rm -rf ./github.com
