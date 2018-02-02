#!/bin/sh

protoc \
  --go_out=plugins=grpc:../proto_test/proto \
  --gostructs_out=import=bitbucket.org/dkfbasel/dev.grpc-tags/proto_test/proto:../proto_test \
  *.proto
