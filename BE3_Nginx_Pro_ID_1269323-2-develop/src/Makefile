.PHONY: all genproto genswag run clean rebuild

all: genproto genswag
	docker compose up -d --build

genproto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/auth.proto
	cp api/auth.pb.go auth/api/
	cp api/auth_grpc.pb.go auth/api/
	cp api/auth.pb.go shop/api/
	cp api/auth_grpc.pb.go shop/api/

genswag:
	cd shop && swag init -d ./cmd,./internal -g main.go --parseInternal --outputTypes go,json,yaml -o ./docs/openapi

run:
	docker compose up -d --build

clean:
	docker compose down -v

rebuild: clean all
