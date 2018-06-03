#!/bin/bash

PARENT_PWD=$(dirname $(pwd))

cd "${PARENT_PWD}/tests"
go test
