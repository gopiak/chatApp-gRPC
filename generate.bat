#!/bin/bash

protoc --proto_path=proto/chatapp --go_out=. --go_opt=module=github.com/webdevgopi/chatApp-gRPC --go-grpc_out=. --go-grpc_opt=module=github.com/webdevgopi/chatApp-gRPC chat_app.proto

protoc --proto_path=proto  --kotlin_out=. proto/chatapp/chat_app.proto