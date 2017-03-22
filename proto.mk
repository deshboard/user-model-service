PROTO_PATH = vendor/github.com/deshboard/apis/deshboard

.PHONY: proto

proto: ## Generate code from protocol buffer
	@mkdir -p apis
	protowrap -I ${PROTO_PATH} ${PROTO_PATH}/user/v1alpha1/user.proto --go_out=plugins=grpc:apis

envcheck::
	$(call executable_check,protoc,protoc)
	$(call executable_check,protowrap,protowrap)
