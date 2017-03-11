PROTO_PATH = vendor/github.com/deshboard/deshboard-proto/proto

.PHONY: proto

proto: ## Generate code from protocol buffer
	@mkdir -p proto
	protowrap -I ${PROTO_PATH} ${PROTO_PATH}/user/user.proto --go_out=plugins=grpc:proto

envcheck::
	$(call executable_check,protoc,protoc)
	$(call executable_check,protowrap,protowrap)
