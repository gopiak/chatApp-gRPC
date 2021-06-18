## Chat implementation gRPC - GoLang - Cassandra

Clone this repository

```bash
cd chatApp-gRPC
go mod tidy
# Starts server - also inserts some mock data into DB
go run main.go  
```

In another terminal ⤵
```bash
go run client/main.go -ID <a userID from cassandraDB>
```
The above command starts a connection between a client of specified ID and the server

The application architecture is as follows:
- When a user connects to the server, a stream will be created (server side streaming) which streams the persisted messages (which are unread) to the connected user
- After recieving unread messages from server, user is prompted to enter an UUID - a unique ID of any user to send and receive msgs realtime
- After entering an UUID, user can send messages to that particular user
- Meanwhile, if someother user sends a msg to the active user, it gets displayed in the console
- Every message is persisted in DB with flag ```is_read True|False```

## Requirements
- Protocol buffer compiler => [Download](https://github.com/protocolbuffers/protobuf/releases)
- Download for windows 32 or 64 bit arch, then unzip into ```C:\```
- In that folder you'll see ```bin``` and ```include``` folders
- **IF WINDOWS** => Add the ```bin``` folder to the path variable - inside Environment variables
- The Go code generator does not produce output for services by default. If you enable the gRPC plugin then code will be generated to support gRPC.
- Download gRPC plugin
- You can run a cassandra image in docker for development. Just add a keySpace ```chat_app```

  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
  ``` 