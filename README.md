# shippy-service-consignment

## Protobuf

Updating the protobuf file requires running something like:

```sh
protoc --micro_out=. --go_out=. proto/consignment/consignment.proto
```

## Build

We utilize a multi-stage build process for deployment. The easiest way to build the service is to use the makefile:

```sh
make build
make run
```