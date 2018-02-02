#!/bin/sh

protoc \
  --gostructs_out=import=github.com/dkfbasel/protobuf/example/domain/proto:../domain \
  --go_out=plugins=grpc:../domain/proto \
  *.proto
