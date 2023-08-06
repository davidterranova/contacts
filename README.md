# Contacts store
This is a very basic contacts store used as a training project to highlight architectures principles
- CLEAN architecture on [main](https://github.com/davidterranova/contacts/tree/main) branch
- CQRS on [cqrs](https://github.com/davidterranova/contacts/tree/cqrs) branch

## WARNING
The store is not thread safe and in-memory

# Highlights
- SOLID/CQRS principles
- Clean architecture
- Clear separation of concerns, layers
- Simple code to read and extend (no complex logic here thought). Quick new member onboarding. 
- Unit tested on key layers without the need to instantiate other layers (possibility to cover more layers, like the API for example)
- DRY
- Self documented

# Run

```
go run main.go server
```

# Dev install

## Protobuff
Install protobuff

https://github.com/protocolbuffers/protobuf

```
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
export PATH="$PATH:$(go env GOPATH)/bin"
```
