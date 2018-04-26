#!/bin/sh

protoc \
  --go_out=. \
  --proto_path=.. \
  empty.proto
