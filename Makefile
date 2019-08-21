build:
	protoc --micro_out=. --go_out=. proto/consignment/consignment.proto
	go build
	GOOS=linux GOARCH=amd64 go build && \
		docker build -t shippy-service-consignment .

run:
	docker run -p 50051:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		shippy-service-consignment
