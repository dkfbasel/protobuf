#!/bin/sh

protoc \
  --go_out=plugins=grpc:../domain \
  *.proto

protoc-go-tags --dir=../domain
