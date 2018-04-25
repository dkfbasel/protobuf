#!/bin/sh

protoc \
  --gostructs_out=../domain \
  --go_out=plugins=grpc:../domain \
  *.proto
