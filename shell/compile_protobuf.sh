#!/bin/bash

PARENT_PWD=$(dirname $(pwd))
PROTO_PWD_ROOT="${PARENT_PWD}/proto"
PROTO_PWD="${PROTO_PWD_ROOT}/protobuf"
cd ${PROTO_PWD}

# 配置一个protoc的软链接到 ${GOPATH}/bin
# 插件及protoc存放路径
BIN_PATH=${GOPATH}/bin
# go get -v github.com/gogo/protobuf/protoc-gen-gogofaster
go install -v github.com/gogo/protobuf/protoc-gen-gogofaster
# go get -v github.com/davyxu/cellnet/protoc-gen-msg
go install -v github.com/davyxu/cellnet/protoc-gen-msg

# 生成协议
${BIN_PATH}/protoc --plugin=protoc-gen-gogofaster=${BIN_PATH}/protoc-gen-gogofaster${EXESUFFIX} --gogofaster_out=. --proto_path="." *.proto
# 生成cellnet 消息注册文件
${BIN_PATH}/protoc --plugin=protoc-gen-msg=${BIN_PATH}/protoc-gen-msg${EXESUFFIX} --msg_out=msgid.go:. --proto_path="." *.proto

find "${PROTO_PWD}" -name "*.pb.go" | xargs -I{} mv {} "${PROTO_PWD_ROOT}/"
mv msgid.go "${PROTO_PWD_ROOT}/"
