#!/bin/sh

protoc \
  --go_out=. \
  --proto_path=.. \
  nullstring.proto

mv ./github.com/dkfbasel/protobuf/types/nullstring/* .
rm -rf ./github.com
