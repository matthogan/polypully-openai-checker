#!/usr/bin/env bash

protoc --go_out=. --go-grpc_out=. protobuf/checker_api.proto
