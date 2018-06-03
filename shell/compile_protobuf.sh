#!/bin/bash

PARENT_PWD=$(dirname $(pwd))

PROTO_PWD="${PARENT_PWD}/protocol"

cd ${PROTO_PWD}"/protobuf"
protoc --go_out="${PROTO_PWD}" *.proto
