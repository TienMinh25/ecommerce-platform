dir_api_gateway := ./internal/api-gateway/migrations
postgres_dsn := postgres://admin:admin@localhost:5432

protoc-compile:
	@protoc \
		--go_out=internal \
		--go_opt=module=github.com/TienMinh25/delivery-system/internal \
		--go-grpc_out=internal \
		--go-grpc_opt=module=github.com/TienMinh25/delivery-system/internal \
		internal/protos/*.proto

migration-create:
	@migrate create -ext sql -dir $(dir_api_gateway) -seq $(name)

migrate-up:
	@migrate -database $(postgres_dsn)/$(dbname)?sslmode=disable -path $(dir_api_gateway) up

migrate-down:
	@migrate -database $(postgres_dsn)/$(dbname)?sslmode=disable -path $(dir_api_gateway) down $(version)

fix-dirty-db:
	@migrate -database $(postgres_dsn)/$(dbname)?sslmode=disable -path $(dir_api_gateway) force $(version)

generate-public-key: generate-private-key
	@openssl rsa -pubout -in jwtRSA256.key -out jwtRSA256.key.pub

generate-private-key:
	@openssl genpkey -algorithm RSA -out jwtRSA256.key

swagger-generate:
	@swag init -g internal/api-gateway/routes/router.go

swagger-format:
	@swag fmt

generate-mock:
	@go generate ./...

tests-run:
	@go test -v -count=1 ./... 2>&1 | grep -v "no test files"

tests-cover:
	@go test -coverprofile=internal/repository/test-cover.out -count=1 -v ./internal/repository/
	@go tool cover -html=internal/repository/test-cover.out

tests-clear:
	@rm internal/repository/test-cover.out

generate-public-key: generate-private-key
	@openssl rsa -pubout -in jwtRSA256.key -out jwtRSA256.key.pub

generate-private-key:
	@openssl genpkey -algorithm RSA -out jwtRSA256.key