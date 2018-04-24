#!/bin/sh

protoc \
  --go_out=. \
  --proto_path=.. \
  timestamp.proto

mv ./github.com/dkfbasel/protobuf/types/timestamp/* .
rm -rf ./github.com
