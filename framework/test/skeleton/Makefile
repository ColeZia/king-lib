GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)

.PHONY: init
# init env
init:
	go get -u github.com/go-kratos/kratos/cmd/kratos/v2
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2

.PHONY: errors
# generate errors code
errors:
	protoc --proto_path=. \
               --proto_path=./third_party \
               --go_out=paths=source_relative:. \
               --go-errors_out=paths=source_relative:. \
               $(API_PROTO_FILES)

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=. \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:. \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=. \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:. \
 	       --go-http_out=paths=source_relative:. \
 	       --go-grpc_out=paths=source_relative:. \
			--openapiv2_out . \
				--openapiv2_opt logtostderr=true \
				--openapiv2_opt json_names_for_fields=false \
	       $(API_PROTO_FILES)

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate
generate:
	go generate ./...

.PHONY: run
# run
run:
	go run cmd/main.go cmd/wire_gen.go --conf conf/app_local.yaml

.PHONY: job-a
# run
job-a:
	go run cmd/jobs/joba/main.go cmd/jobs/joba/wire_gen.go --conf conf/app_local.yaml

.PHONY: job-b
# run
job-b:
	go run cmd/jobs/jobb/main.go cmd/jobs/jobb/wire_gen.go --conf conf/app_local.yaml


.PHONY: ent-demo-gen
ent-demo-gen:
	go generate ./internal/data/ent/demo/

#数据迁移
.PHONY: migrate
# run
migrate:
	go run cmd/boss-test/main.go cmd//boss-test/wire_gen.go migrate --conf conf/app_local.yaml

.PHONY: migrate-sync
# run
migrate-sync:
	go run cmd/boss-test/main.go cmd//boss-test/wire_gen.go migrate --conf conf/app_local.yaml -m sync

.PHONY: update-pb
update-pb:
	go get gl.king.im/king-lib/protobuf

.PHONY: update-fw
update-fw:
	go get gl.king.im/king-lib/framework
	
.PHONY: update-pkg
update-pkg:
	go get gl.king.im/king-lib/pkg

.PHONY: all
# generate all
all:
	make api;
	make errors;
	make config;
	make generate;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
