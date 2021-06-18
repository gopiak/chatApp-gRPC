#!/bin/bash

protoc --proto_path=proto/chatapp --go_out=. --go_opt=module=api/microservices/chatApp-gRPC --go-grpc_out=. --go-grpc_opt=module=api/microservices/chatApp-gRPC chat_app.proto

protoc --proto_path=proto  --kotlin_out=. proto/chatapp/chat_app.proto